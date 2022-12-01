// Package assert provides test assertion helpers.
package assert

import (
	"io"
	"path/filepath"
	"testing"

	"oss.terrastruct.com/util-go/xjson"

	"oss.terrastruct.com/util-go/diff"
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
	err := diff.TestdataJSON(filepath.Join("testdata", tb.Name()), got)
	Success(tb, err)
}

func Testdata(tb testing.TB, ext string, got []byte) {
	err := diff.Testdata(filepath.Join("testdata", tb.Name()), ext, got)
	Success(tb, err)
}

func Close(t *testing.T, c io.Closer) {
	err := c.Close()
	if err != nil {
		t.Fatalf("failed to close %T: %v", c, err)
	}
}
