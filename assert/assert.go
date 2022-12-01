// Package assert provides test assertion helpers.
package assert

import (
	"path/filepath"
	"testing"

	"oss.terrastruct.com/diff"
	"oss.terrastruct.com/xjson"
)

func Success(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func JSON(t testing.TB, exp string, got interface{}) {
	t.Helper()
	String(t, exp, xjson.MarshalIndent(got))
}

func String(t testing.TB, exp, got string) {
	t.Helper()
	if exp == got {
		return
	}
	diff, err := diff.Strings(exp, got)
	if err != nil {
		t.Fatal(err)
	}
	if diff != "" {
		t.Fatalf("\n%s", diff)
	}
}

func Testdata(t testing.TB, got interface{}) {
	err := diff.Testdata(filepath.Join("testdata", t.Name()), got)
	Success(t, err)
}
