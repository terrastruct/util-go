package cmdlog_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/creack/pty"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/cmdlog"
	"oss.terrastruct.com/util-go/xos"
	"oss.terrastruct.com/util-go/xtesting"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	var cases = []xtesting.Case{
		{
			Name: "COLOR=1",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				b := &bytes.Buffer{}
				env.Setenv("COLOR", "1")
				l := cmdlog.New(env, b)

				testLogger(l)

				t.Log(b.String())
				assert.TestdataJSON(t, b.String())
			},
		},
		{
			Name: "COLOR=",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				b := &bytes.Buffer{}
				l := cmdlog.New(env, b)

				testLogger(l)

				t.Log(b.String())
				assert.TestdataJSON(t, b.String())
			},
		},
		{
			Name: "tty",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
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
				assert.TestdataJSON(t, string(out))
			},
		},
		{
			Name: "testing.TB",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
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
			Name: "WithPrefix",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
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
				assert.TestdataJSON(t, b.String())
			},
		},
		{
			Name: "multiline",
			Run: func(t *testing.T, ctx context.Context, env *xos.Env) {
				b := &bytes.Buffer{}
				env.Setenv("COLOR", "1")
				l := cmdlog.New(env, b)

				l.NoLevel.Print("")
				l.SetTS(true)
				l.NoLevel.Print("")
				l.SetTS(false)

				l2 := l.WithCCPrefix("lochness")
				l2 = l2.WithCCPrefix("imgbundler")
				l2 = l2.WithCCPrefix("cache")

				l2.Warn.Print(``)
				l2.Warn.Print("\n\n\n")
				l2.SetTS(true)
				l2.Warn.Printf(`yes %d
yes %d`, 3, 4)

				t.Log(b.String())
				assert.TestdataJSON(t, b.String())
			},
		},
	}

	xtesting.RunCases(t, cases)
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
	l.Error.Println(`Good day to let down old friends who need help.
I believe in getting into hot water; it keeps you clean.`)
}

type fakeTB struct {
	testing.TB
	logf func(string, ...interface{})
}

func (ftb *fakeTB) Logf(f string, v ...interface{}) {
	ftb.TB.Helper()
	ftb.logf(f, v...)
}
