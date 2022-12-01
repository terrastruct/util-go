package xdefer_test

import (
	"fmt"
	"runtime"
	"testing"

	"golang.org/x/xerrors"

	"oss.terrastruct.com/utils-go/assert"
	"oss.terrastruct.com/utils-go/xdefer"
)

func TestErrorf(t *testing.T) {
	t.Parallel()

	err := func() (err error) {
		defer xdefer.Errorf(&err, "second wrap %#v", []int{99, 3})

		err = xerrors.New("ola amigo")
		if err != nil {
			// This is the first line that should be reported on xdefer.
			return xerrors.Errorf("first wrap: %w", err)
		}

		return nil
	}()

	_, fp, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	exp := fmt.Sprintf(`second wrap []int{99, 3}:
  - first wrap:
    oss.terrastruct.com/utils-go/xdefer_test.TestErrorf.func1
        %v:23
  - ola amigo:
    oss.terrastruct.com/utils-go/xdefer_test.TestErrorf.func1
        %[1]v:20`,
		fp,
	)

	got := fmt.Sprintf("%+v", err)
	assert.String(t, exp, got)
}

func TestEmptyErrorf(t *testing.T) {
	t.Parallel()

	err := func() (err error) {
		defer xdefer.Errorf(&err, "")

		err = xerrors.New("ola amigo")
		if err != nil {
			return err
		}

		return nil
	}()

	_, fp, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	exp := fmt.Sprintf(`oss.terrastruct.com/utils-go/xdefer_test.TestEmptyErrorf.func1
        %v:55
  - ola amigo:
    oss.terrastruct.com/utils-go/xdefer_test.TestEmptyErrorf.func1
        %[1]v:53`,
		fp,
	)

	got := fmt.Sprintf("%+v", err)
	assert.String(t, exp, got)

	exp = err.Error()
	got = "oss.terrastruct.com/utils-go/xdefer_test.TestEmptyErrorf.func1: ola amigo"
	assert.String(t, exp, got)
}
