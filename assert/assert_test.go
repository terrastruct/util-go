package assert_test

import (
	"testing"

	"oss.terrastruct.com/util-go/xjson"

	"oss.terrastruct.com/util-go/assert"
)

func TestStringJSON(t *testing.T) {
	t.Parallel()

	gen := func() (m1, m2 map[string]interface{}) {
		m1 = map[string]interface{}{
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

		m2 = map[string]interface{}{
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

		return m1, m2
	}

	t.Run("equal", func(t *testing.T) {
		t.Parallel()

		m1, m2 := gen()
		assert.StringJSON(t, xjson.MarshalIndent(m1), m2)
	})

	t.Run("diff", func(t *testing.T) {
		t.Parallel()

		m1, m2 := gen()
		m2["five"].(map[string]interface{})["five"].(map[string]interface{})["no"] = "ys"

		fataledWithDiff := false
		ftb := &fakeTB{
			TB: t,
			fatalf: func(f string, v ...interface{}) {
				t.Helper()
				if len(v) == 1 {
					t.Logf(f, v...)
					fataledWithDiff = true
					return
				}

				t.Fatalf(f, v...)
			}}

		defer func() {
			if t.Failed() || !fataledWithDiff {
				t.Error("expected assert.StringJSON to fatal with correct diff")
			}
		}()
		assert.StringJSON(ftb, xjson.MarshalIndent(m1), m2)
	})
}

type fakeTB struct {
	fatalf func(string, ...interface{})
	testing.TB
}

func (ftb *fakeTB) Fatalf(f string, v ...interface{}) {
	ftb.TB.Helper()
	ftb.fatalf(f, v...)
}
