package xmain_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
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
			name: "hello_world",
			run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				env.Setenv("HELLO_FLAG", "world")
				ts := &xmain.TestingState{
					Run: func(ctx context.Context, ms *xmain.State) error {
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
					},
					Env:  env,
					Args: []string{"hello"},
				}
				stdout := ts.StdoutPipe()

				ts.Start(t, ctx)
				defer ts.Cleanup(t)

				s, err := io.ReadAll(stdout)
				assert.Success(t, err)

				err = ts.Wait(ctx)
				assert.Success(t, err)

				assert.Equal(t, "world", string(s))
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
