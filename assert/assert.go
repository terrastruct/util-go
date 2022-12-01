// Package assert provides test assertion helpers.
package assert

import (
	"testing"

	"oss.terrastruct.com/xjson"
	"oss.terrastruct.com/diff"
)

func Success(t *testing.T, err error) {
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
