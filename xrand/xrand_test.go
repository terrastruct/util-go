package xrand_test

import (
	"testing"

	"oss.terrastruct.com/utils-go/xrand"
)

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("invalids", func(t *testing.T) {
		t.Parallel()

		s := xrand.String(0, nil)
		if s != "" {
			t.Fatalf("expected empty string: %q", s)
		}
		s = xrand.String(-1, nil)
		if s != "" {
			t.Fatalf("expected empty string: %q", s)
		}
	})

	s := xrand.String(20, nil)
	t.Logf("%d %s", len(s), s)
}
