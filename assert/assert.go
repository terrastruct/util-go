// Package assert provides test assertion helpers.
package assert

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"oss.terrastruct.com/util-go/diff"
	"oss.terrastruct.com/util-go/xjson"
)

func Success(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("unexpected error: %v", err)
	}
}

func Error(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		tb.Fatal("expected error")
	}
}

func ErrorString(tb testing.TB, err error, msg string) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("expected error containing %q", msg)
	}
	String(tb, msg, err.Error())
}

func StringJSON(tb testing.TB, exp string, got interface{}) {
	tb.Helper()
	String(tb, exp, string(xjson.Marshal(got)))
}

func String(tb testing.TB, exp, got string) {
	tb.Helper()
	diff, err := diff.Strings(exp, got)
	Success(tb, err)
	if diff != "" {
		tb.Fatalf("\n%s", diff)
	}
}

func JSON(tb testing.TB, exp, got interface{}) {
	tb.Helper()
	diff, err := diff.JSON(exp, got)
	Success(tb, err)
	if diff != "" {
		tb.Fatalf("\n%s", diff)
	}
}

func Runes(tb testing.TB, exp, got string) {
	tb.Helper()
	err := diff.Runes(exp, got)
	Success(tb, err)
}

func TestdataJSON(tb testing.TB, got interface{}) {
	tb.Helper()
	err := diff.TestdataJSON(filepath.Join("testdata", tb.Name()), got)
	Success(tb, err)
}

func Testdata(tb testing.TB, ext string, got []byte) {
	tb.Helper()
	err := diff.Testdata(filepath.Join("testdata", tb.Name()), ext, got)
	Success(tb, err)
}

func TestdataDir(tb testing.TB, dir string) {
	tb.Helper()
	testdataDir(tb, filepath.Join("testdata", tb.Name()), dir)
	if tb.Failed() {
		tb.FailNow()
	}
}

func testdataDir(tb testing.TB, testName, dir string) {
	tb.Helper()

	ea, err := os.ReadDir(dir)
	Success(tb, err)

	for _, e := range ea {
		if e.IsDir() {
			testdataDir(tb, filepath.Join(testName, e.Name()), filepath.Join(dir, e.Name()))
			continue
		}
		ext := filepath.Ext(e.Name())
		name := strings.TrimSuffix(e.Name(), ext)
		got := ReadFile(tb, filepath.Join(dir, e.Name()))
		err := diff.Testdata(filepath.Join(testName, name), ext, got)
		if err != nil {
			tb.Error(err)
		}
	}
}

func Close(tb testing.TB, c io.Closer) {
	tb.Helper()
	err := c.Close()
	if err != nil {
		tb.Fatalf("failed to close %T: %v", c, err)
	}
}

func Equal(tb testing.TB, exp, got interface{}) {
	tb.Helper()
	if exp == got {
		return
	}
	exps, ok := exp.(string)
	if ok {
		gots, ok := got.(string)
		if ok {
			String(tb, exps, gots)
			return
		}
	}
	tb.Fatalf("expected %#v but got %#v", exp, got)
}

func NotEqual(tb testing.TB, v1, v2 interface{}) {
	tb.Helper()
	if v1 != v2 {
		return
	}
	tb.Fatalf("did not expect %#v", v2)
}

func True(tb testing.TB, v bool) {
	tb.Helper()
	Equal(tb, true, v)
}

func False(tb testing.TB, v bool) {
	tb.Helper()
	Equal(tb, false, v)
}

func TempDir(tb testing.TB) (dir string, cleanup func()) {
	tb.Helper()

	dir, err := os.MkdirTemp("", "util-go.assert.TempDir")
	Success(tb, err)
	return dir, func() {
		tb.Helper()
		err = os.RemoveAll(dir)
		Success(tb, err)
	}
}

func WriteFile(tb testing.TB, fp string, data []byte, perms os.FileMode) {
	tb.Helper()

	err := os.WriteFile(fp, data, perms)
	Success(tb, err)
}

func ReadFile(tb testing.TB, fp string) (data []byte) {
	tb.Helper()

	data, err := os.ReadFile(fp)
	Success(tb, err)
	return data
}

func Remove(tb testing.TB, fp string) {
	tb.Helper()

	err := os.Remove(fp)
	Success(tb, err)
}

func RemoveAll(tb testing.TB, fp string) {
	tb.Helper()

	err := os.RemoveAll(fp)
	Success(tb, err)
}
