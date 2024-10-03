// package diff contains diff generation helpers, particularly useful for tests.
//
// - Strings
// - Files
// - Runes
// - JSON
// - Testdata
// - TestdataJSON
package diff

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/multierr"

	"oss.terrastruct.com/util-go/xdefer"
	"oss.terrastruct.com/util-go/xjson"
)

// Strings diffs exp with got in a git style diff.
//
// The git style diff header will contain real paths to exp and got
// on the file system so that you can easily inspect them.
//
// This behavior is particularly useful for when you need to update
// a test with the new got. You can just copy and paste from the got
// file in the diff header.
//
// It uses Files under the hood.
func Strings(exp, got string) (ds string, err error) {
	defer xdefer.Errorf(&err, "failed to diff text")

	if exp == got {
		return "", nil
	}

	d, err := ioutil.TempDir("", "ts_d2_diff")
	if err != nil {
		return "", err
	}

	expPath := filepath.Join(d, "exp")
	gotPath := filepath.Join(d, "got")

	err = ioutil.WriteFile(expPath, []byte(exp), 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(gotPath, []byte(got), 0644)
	if err != nil {
		return "", err
	}

	return Files(expPath, gotPath)
}

// Files diffs expPath with gotPath and prints a git style diff header.
//
// It uses git under the hood.
func Files(expPath, gotPath string) (ds string, err error) {
	defer xdefer.Errorf(&err, "failed to diff files")

	_, err = os.Stat(expPath)
	if os.IsNotExist(err) {
		expPath = "/dev/null"
	}
	_, err = os.Stat(gotPath)
	if os.IsNotExist(err) {
		gotPath = "/dev/null"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-c", "diff.color=always", "diff",
		// Use the best diff-algorithm and highlight trailing whitespace.
		"--diff-algorithm=histogram",
		"--ws-error-highlight=all",
		"--no-index",
		expPath, gotPath)
	cmd.Env = append(cmd.Env, "GIT_CONFIG_NOSYSTEM=1", "HOME=")

	diffBytes, err := cmd.CombinedOutput()
	var ee *exec.ExitError
	if err != nil && !errors.As(err, &ee) {
		return "", fmt.Errorf("git diff failed: out=%q: %w", diffBytes, err)
	}
	ds = string(diffBytes)

	// Strips the diff header before ---
	//
	// diff --git a/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/exp b/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/got
	// index d48c704b..dbe709e6 100644
	// --- a/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/exp
	// +++ b/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/got
	// @@ -1,5 +1,5 @@
	//
	// becomes:
	//
	// --- a/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/exp
	// +++ b/var/folders/tf/yp9nqwbx4g5djqjxms03cvx80000gn/T/d2parser_test916758829/got
	// @@ -1,5 +1,5 @@
	i := strings.Index(ds, "index")
	if i > -1 {
		j := strings.IndexByte(ds[i:], '\n')
		if j > -1 {
			ds = ds[i+j+1:]
		}
	}
	return strings.TrimSpace(ds), nil
}

// Runes is like Strings but formats exp and got with each unicode codepoint on a separate
// line and generates a diff of that. It's useful for autogenerated UTF-8 with
// xrand.String as Strings won't generate a coherent diff with undisplayable characters.
func Runes(exp, got string) error {
	if exp == got {
		return nil
	}
	expRunes := formatRunes(exp)
	gotRunes := formatRunes(got)
	ds, err := Strings(expRunes, gotRunes)
	if err != nil {
		return err
	}
	if ds != "" {
		return errors.New(ds)
	}
	return nil
}

func formatRunes(s string) string {
	return strings.Join(strings.Split(fmt.Sprintf("%#v", []rune(s)), ", "), "\n")
}

// TestdataJSON is for when you have JSON that is too large to easily keep embedded by the
// tests in _test.go files. As well, it makes the acceptance of large changes trivial
// unlike say fs/embed.
//
// TestdataJSON encodes got as JSON and diffs it against the stored json in path.exp.json.
// The got JSON is stored in path.got.json. If the diff is empty, it returns nil.
//
// Otherwise it returns an error containing the diff.
//
// In order to accept changes path.got.json has to become path.exp.json. You can use
// ./ci/testdata/accept.sh to rename all non stale path.got.json files to path.exp.json.
//
// You can scope it to a single test or folder, see ./ci/testdata/accept.sh --help
//
// Also see ./ci/testdata/clean.sh --help for cleaning the repository of all
// path.got.json and path.exp.json files.
//
// You can also use $TESTDATA_ACCEPT=1 to update all path.exp.json files on the fly.
// This is useful when you're regenerating the repository's testdata. You can't easily
// use the accept script without rerunning go test multiple times as go test will return
// after too many test failures and will not continue until they are fixed.
//
// You'll want to use -count=1 to disable go test's result caching if you do use
// $TESTDATA_ACCEPT.
//
// TestdataJSON will automatically create nonexistent directories in path.
//
// Here's an example that you can play with to better understand the behaviour:
//
//	err = diff.TestdataJSON(filepath.Join("testdata", t.Name()), "change me")
//	if err != nil {
//		t.Fatal(err)
//	}
//
// Normally you want to use t.Name() as path for clarity but you can pass in any string.
// e.g. a single test could persist two json objects into testdata with:
//
//	err = diff.TestdataJSON(filepath.Join("testdata", t.Name(), "1"), "change me 1")
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = diff.TestdataJSON(filepath.Join("testdata", t.Name(), "2"), "change me 2")
//	if err != nil {
//		t.Fatal(err)
//	}
//
// These would persist in testdata/${t.Name()}/1.exp.json and testdata/${t.Name()}/2.exp.json
//
// It uses Files under the hood.
//
// note: testdata is the canonical Go directory for such persistent test only files.
//
//	It is unfortunately poorly documented. See https://pkg.go.dev/cmd/go/internal/test
//	So normally you'd want path to be filepath.Join("testdata", t.Name()).
//	This is also the reason this function is named "TestdataJSON".
func TestdataJSON(path string, got interface{}) error {
	gotb := xjson.Marshal(got)
	gotb = append(gotb, '\n')
	return Testdata(path, ".json", gotb)
}

// ext includes period like path.Ext()
func Testdata(path, ext string, got []byte) error {
	expPath := fmt.Sprintf("%s.exp%s", path, ext)
	gotPath := fmt.Sprintf("%s.got%s", path, ext)

	err := os.MkdirAll(filepath.Dir(gotPath), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(gotPath, []byte(got), 0600)
	if err != nil {
		return err
	}

	ds, err := Files(expPath, gotPath)
	if err != nil {
		return err
	}

	if ds != "" {
		if os.Getenv("TESTDATA_ACCEPT") != "" || os.Getenv("TA") != "" {
			return os.Rename(gotPath, expPath)
		}
		if os.Getenv("NO_DIFF") != "" || os.Getenv("ND") != "" {
			ds = "diff hidden with $NO_DIFF=1 or $ND=1"
		}
		return fmt.Errorf("diff (rerun with $TESTDATA_ACCEPT=1 or $TA=1 to accept):\n%s", ds)
	}
	return os.Remove(gotPath)
}

func JSON(exp, got interface{}) (string, error) {
	return Strings(string(xjson.Marshal(exp)), string(xjson.Marshal(got)))
}

func TestdataDir(testName, dir string) (err error) {
	defer xdefer.Errorf(&err, "failed to commit testdata dir %v", dir)
	testdataDir(&err, testName, dir)
	return err
}

func testdataDir(errs *error, testName, dir string) {
	ea, err := os.ReadDir(dir)
	if err != nil {
		*errs = multierr.Combine(*errs, err)
		return
	}

	for _, e := range ea {
		if e.IsDir() {
			testdataDir(errs, filepath.Join(testName, e.Name()), filepath.Join(dir, e.Name()))
		} else {
			ext := filepath.Ext(e.Name())
			name := strings.TrimSuffix(e.Name(), ext)
			got, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				*errs = multierr.Combine(*errs, err)
				continue
			}
			err = Testdata(filepath.Join(testName, name), ext, got)
			if err != nil {
				*errs = multierr.Combine(*errs, err)
			}
		}
	}
}
