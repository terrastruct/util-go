package diff_test

import (
	"os"
	"path/filepath"
	"testing"

	"oss.terrastruct.com/utils-go/assert"
	"oss.terrastruct.com/utils-go/diff"
)

//lint:file-ignore ST1018 ignore staticcheck string literal with Unicode control characters

func TestTestData(t *testing.T) {
	t.Run("TESTDATA_ACCEPT", testTestDataAccept)

	os.Unsetenv("TESTDATA_ACCEPT")

	m1 := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five": map[string]interface{}{
			"yes": "yes",
			"no":  "yes",
			"five": map[string]interface{}{
				"yes": "no",
				"no":  "yes",
			},
		},
	}

	err := os.Remove("testdata/TestTestData.exp.json")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("unexpected error: %v", err)
	}

	err = diff.Testdata(filepath.Join("testdata", t.Name()), m1)
	assert.Error(t, err)
	exp := `diff (rerun with $TESTDATA_ACCEPT=1 to accept):
[1m--- /dev/null[m
[1m+++ b/testdata/TestTestData.got.json[m
[36m@@ -0,0 +1,14 @@[m
[32m+[m[32m{[m
[32m+[m[32m  "five": {[m
[32m+[m[32m    "five": {[m
[32m+[m[32m      "no": "yes",[m
[32m+[m[32m      "yes": "no"[m
[32m+[m[32m    },[m
[32m+[m[32m    "no": "yes",[m
[32m+[m[32m    "yes": "yes"[m
[32m+[m[32m  },[m
[32m+[m[32m  "four": 4,[m
[32m+[m[32m  "one": 1,[m
[32m+[m[32m  "three": 3,[m
[32m+[m[32m  "two": 2[m
[32m+[m[32m}[m`
	got := err.Error()
	ds, err := diff.Strings(exp, got)
	if err != nil {
		t.Fatalf("unable to generate exp diff: %v", err)
	}
	if ds != "" {
		t.Fatalf("expected no diff:\n%s", ds)
	}
	err = diff.Runes(exp, got)
	if err != nil {
		t.Fatalf("expected no rune diff: %v", err)
	}

	err = os.Rename("testdata/TestTestData.got.json", "testdata/TestTestData.exp.json")
	assert.Success(t, err)

	m1["five"].(map[string]interface{})["five"].(map[string]interface{})["no"] = "ys"

	err = diff.Testdata(filepath.Join("testdata", t.Name()), m1)
	if err == nil {
		t.Fatalf("expected err: %#v", err)
	}
	exp = `diff (rerun with $TESTDATA_ACCEPT=1 to accept):
[1m--- a/testdata/TestTestData.exp.json[m
[1m+++ b/testdata/TestTestData.got.json[m
[36m@@ -1,7 +1,7 @@[m
 [m{[m
 [m  "five": {[m
 [m    "five": {[m
[31m-[m[31m      "no": "yes",[m
[32m+[m[32m      "no": "ys",[m
 [m      "yes": "no"[m
 [m    },[m
 [m    "no": "yes",[m`
	got = err.Error()
	ds, err = diff.Strings(exp, got)
	assert.Success(t, err)
	if ds != "" {
		t.Fatalf("expected no diff:\n%s", ds)
	}

	exp += "a"
	ds, err = diff.Strings(exp, got)
	assert.Success(t, err)
	if ds == "" {
		t.Fatalf("expected incorrect diff:\n%s", ds)
	}
	err = diff.Runes(exp, got)
	assert.Error(t, err)
}

func testTestDataAccept(t *testing.T) {
	m1 := map[string]interface{}{
		"one": 1,
	}

	os.Setenv("TESTDATA_ACCEPT", "1")
	err := diff.Testdata(filepath.Join("testdata", t.Name()), m1)
	assert.Success(t, err)

	m1["one"] = 2

	os.Setenv("TESTDATA_ACCEPT", "")
	err = diff.Testdata(filepath.Join("testdata", t.Name()), m1)
	assert.Error(t, err)
	exp := `diff (rerun with $TESTDATA_ACCEPT=1 to accept):
[1m--- a/testdata/TestTestData/TESTDATA_ACCEPT.exp.json[m
[1m+++ b/testdata/TestTestData/TESTDATA_ACCEPT.got.json[m
[36m@@ -1,3 +1,3 @@[m
 [m{[m
[31m-[m[31m  "one": 1[m
[32m+[m[32m  "one": 2[m
 [m}[m`
	ds, err := diff.Strings(exp, err.Error())
	assert.Success(t, err)
	if ds != "" {
		t.Fatalf("expected no diff:\n%s", ds)
	}

	err = os.Remove(filepath.Join("testdata", t.Name()) + ".got.json")
	assert.Success(t, err)
}
