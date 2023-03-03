package xmain_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/xmain"
	"oss.terrastruct.com/util-go/xos"
)

func TestTesting(t *testing.T) {
	t.Parallel()

	tca := []struct {
		name string
		run  func(t *testing.T, ctx context.Context, env *xos.Env)
	}{
		{
			name: "base",
			run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				ts := &xmain.TestState{
					Run:  helloWorldRun,
					Env:  env,
					Args: []string{"helloWorldRun"},
				}

				ts.Start(t, ctx)
				defer ts.Cleanup(t)

				err := ts.Wait(ctx)
				assert.ErrorString(t, err, `failed to wait testing xmain: helloWorldRun: bad usage: $HELLO_FLAG or -flag missing`)
			},
		},
		{
			name: "help",
			run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				stdout := &strings.Builder{}
				ts := &xmain.TestState{
					Run:    helloWorldRun,
					Env:    env,
					Args:   []string{"helloWorldRun", "-help"},
					Stdout: stdout,
				}

				ts.Start(t, ctx)
				defer ts.Cleanup(t)

				err := ts.Wait(ctx)
				assert.Success(t, err)

				assert.Equal(t, `Usage:
helloWorldRun [-flag=val]

helloWorldRun prints the value of -flag to stdout. $HELLO_FLAG is equivalent to -flag.
`, stdout.String())
			},
		},
		{
			name: "envPriority",
			run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				env.Setenv("HELLO_FLAG", "world")
				stdout := &strings.Builder{}
				ts := &xmain.TestState{
					Run:    helloWorldRun,
					Env:    env,
					Args:   []string{"helloWorldRun", "hello"},
					Stdout: stdout,
				}

				ts.Start(t, ctx)
				defer ts.Cleanup(t)

				err := ts.Wait(ctx)
				assert.Success(t, err)

				assert.Equal(t, "world", stdout.String())
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tca {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			tc.run(t, ctx, xos.NewEnv(nil))
		})
	}
}

func helloWorldRun(ctx context.Context, ms *xmain.State) error {
	flag := ms.Opts.String("HELLO_FLAG", "flag", "f", "", "")
	err := ms.Opts.Flags.Parse(ms.Opts.Args)
	if errors.Is(err, pflag.ErrHelp) {
		fmt.Fprintf(ms.Stdout, `Usage:
%[1]s [-flag=val]

%[1]s prints the value of -flag to stdout. $HELLO_FLAG is equivalent to -flag.
`, filepath.Base(ms.Name))
		return nil
	}
	if *flag == "" {
		return xmain.UsageErrorf("$HELLO_FLAG or -flag missing")
	}
	_, err = io.WriteString(ms.Stdout, *flag)
	return err
}
