// Package assert provides test assertion helpers.
package assert

import "testing"

func Success(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
