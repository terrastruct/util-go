package xmain_test

import (
	"os"
	"testing"

	"oss.terrastruct.com/util-go/diff"
	"oss.terrastruct.com/util-go/xmain"
	"oss.terrastruct.com/util-go/xos"
)

func TestOpts(t *testing.T) {
	// TODO
	t.Skip()
	ms := &xmain.State{
		Name: "test",
		Env:  xos.NewEnv(os.Environ()),
	}
	ms.Opts = xmain.NewOpts(ms.Env, nil, nil)

	// Test that it works when strings are longer than env vars
	ms.Opts.String("D2_LAYOUT", "layout", "l", "dagre", "this is the layout")
	ms.Opts.String("D2_THEME", "theme", "t", "neutral_default", "theme")
	ms.Opts.String("", "something-very-very-long-lalalal", "", "test", "this is a long string")

	got := ms.Opts.Defaults()
	ds, err := diff.Strings("", got)
	if err != nil {
		t.Fatal(err)
	}
	if ds != "" {
		t.Fatalf("got != exp:\n%s", ds)
	}
}
