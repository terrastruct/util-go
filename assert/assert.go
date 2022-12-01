// Package assert provides test assertion helpers.
package assert

import (
	"io"
	"path/filepath"
	"testing"

	"oss.terrastruct.com/utils-go/xjson"

	"oss.terrastruct.com/utils-go/diff"
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

func ErrorContains(tb testing.TB, err error, msg string) {
	tb.Helper()
	if err == nil {
		tb.Fatalf("expected error containing %q", msg)
	}
	String(tb, msg, err.Error())
}

func JSON(tb testing.TB, exp, got interface{}) {
	tb.Helper()
	String(tb, xjson.MarshalIndent(exp), xjson.MarshalIndent(got))
}

func StringJSON(tb testing.TB, exp string, got interface{}) {
	tb.Helper()
	String(tb, exp, xjson.MarshalIndent(got))
}

func String(tb testing.TB, exp, got string) {
	tb.Helper()
	diff, err := diff.Strings(exp, got)
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

func Testdata(tb testing.TB, got interface{}) {
	err := diff.Testdata(filepath.Join("testdata", tb.Name()), got)
	Success(tb, err)
}

func Close(t *testing.T, c io.Closer) {
	err := c.Close()
	if err != nil {
		t.Fatalf("failed to close %T: %v", c, err)
	}
}
