package cmdlog_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/creack/pty"

	"oss.terrastruct.com/utils-go/assert"
	"oss.terrastruct.com/utils-go/cmdlog"
	"oss.terrastruct.com/utils-go/xos"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		name string

		run func(*testing.T, *xos.Env)
	}{
		{
			name: "COLOR=1",
			run: func(t *testing.T, env *xos.Env) {
				b := &bytes.Buffer{}
				env.Setenv("COLOR", "1")
				l := cmdlog.New(env, b)

				testLogger(l)

				t.Log(b.String())
				assert.Testdata(t, b.String())
			},
		},
		{
			name: "COLOR=",
			run: func(t *testing.T, env *xos.Env) {
				b := &bytes.Buffer{}
				l := cmdlog.New(env, b)

				testLogger(l)

				t.Log(b.String())
				assert.Testdata(t, b.String())
			},
		},
		{
			name: "tty",
			run: func(t *testing.T, env *xos.Env) {
				ptmx, tty, err := pty.Open()
				if err != nil {
					t.Fatalf("failed to open pty: %v", err)
				}
				defer assert.Close(t, ptmx)
				defer assert.Close(t, tty)

				l := cmdlog.New(env, tty)
				testLogger(l)

				timer := time.AfterFunc(time.Second*5, func() {
					// For some reason ptmx.SetDeadline() does not work.
					// tty has to be closed for a read on ptmx to unblock.
					tty.Close()
					t.Error("read took too long, update expLen")
				})
				defer timer.Stop()
				// If the expected output changes, increase this to 9999, rerun and then update
				// for new length.
				const expLen = 415
				out, err := io.ReadAll(io.LimitReader(ptmx, expLen))
				if err != nil {
					t.Fatalf("failed to read log output: %v", err)
				}
				t.Log(len(out), string(out))
				assert.Testdata(t, string(out))
			},
		},
		{
			name: "testing.TB",
			run: func(t *testing.T, env *xos.Env) {
				ft := &fakeTB{
					TB: t,
					logf: func(f string, v ...interface{}) {
						t.Helper()
						assert.String(t, "info: what's up\n", fmt.Sprintf(f, v...))
					},
				}

				env.Setenv("COLOR", "0")
				l := cmdlog.NewTB(env, ft)
				l.Info.Printf("what's up")
			},
		},
		{
			name: "WithPrefix",
			run: func(t *testing.T, env *xos.Env) {
				b := &bytes.Buffer{}
				env.Setenv("COLOR", "1")
				l := cmdlog.New(env, b)

				l2 := l.WithCCPrefix("lochness")
				if l2 == l {
					t.Fatalf("expected l and l2 to be different loggers")
				}
				l2 = l2.WithCCPrefix("imgbundler")
				l2 = l2.WithCCPrefix("cache")

				testLogger(l)
				testLogger(l2)

				t.Log(b.String())
				assert.Testdata(t, b.String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			env := xos.NewEnv(nil)
			tc.run(t, env)
		})
	}
}

func testLogger(l *cmdlog.Logger) {
	l.NoLevel.Println("Somehow, the world always affects you more than you affect it.")

	l.SetDebug(true)
	l.Debug.Println("Man is a rational animal who always loses his temper when he is called upon.")

	l.SetDebug(false)
	l.Debug.Println("You can never trust a woman; she may be true to you.")

	l.SetTS(true)
	l.Success.Println("An alcoholic is someone you don't like who drinks as much as you do.")
	l.Info.Println("There once was this swami who lived above a delicatessan.")

	l.SetTSFormat(time.UnixDate)
	l.Warn.Println("Telephone books are like dictionaries -- if you know the answer before.")

	l.SetTS(false)
	l.Error.Println("Nothing can be done in one trip.")
}

type fakeTB struct {
	testing.TB
	logf func(string, ...interface{})
}

func (ftb *fakeTB) Logf(f string, v ...interface{}) {
	ftb.TB.Helper()
	ftb.logf(f, v...)
}
