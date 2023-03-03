package xmain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/cmdlog"
	"oss.terrastruct.com/util-go/xcontext"
	"oss.terrastruct.com/util-go/xdefer"
	"oss.terrastruct.com/util-go/xos"
)

type TestingState struct {
	Run  func(context.Context, *State) error
	Env  *xos.Env
	Args []string
	Dir  string
	FS   fs.FS

	Stdin  io.Reader
	Stdout io.WriteCloser
	Stderr io.WriteCloser

	mu      *xcontext.Mutex
	ms      *State
	sigs    chan os.Signal
	done    chan error
	doneErr *error
}

func (ts *TestingState) StdinPipe() (pw io.WriteCloser) {
	ts.Stdin, pw = io.Pipe()
	return pw
}

func (ts *TestingState) StdoutPipe() (pr io.Reader) {
	pr, ts.Stdout = io.Pipe()
	return pr
}

func (ts *TestingState) StderrPipe() (pr io.Reader) {
	pr, ts.Stderr = io.Pipe()
	return pr
}

func (ts *TestingState) Start(tb testing.TB, ctx context.Context) {
	tb.Helper()

	if ts.mu != nil {
		tb.Fatal("xmain.TestingState.Start cannot be called twice")
	}
	if ts.Env == nil {
		ts.Env = xos.NewEnv(nil)
	}

	ts.mu = xcontext.NewMutex()
	ts.sigs = make(chan os.Signal, 1)
	ts.done = make(chan error, 1)

	name := ""
	args := []string(nil)
	if len(args) > 0 {
		name = os.Args[0]
		args = os.Args[1:]
	}
	log := cmdlog.NewTB(ts.Env, tb)
	ts.ms = &State{
		Name: name,

		Log:  log,
		Env:  ts.Env,
		Opts: NewOpts(ts.Env, log, args),
		Dir:  ts.Dir,
		FS:   ts.FS,
	}

	ts.ms.Stdin = ts.Stdin
	if ts.Stdin == nil {
		ts.ms.Stdin = io.LimitReader(nil, 0)
	}
	ts.ms.Stdout = ts.Stdout
	if ts.Stdout == nil {
		ts.ms.Stdout = nopWriterCloser{io.Discard}
	}
	ts.ms.Stderr = ts.Stderr
	if ts.Stderr == nil {
		ts.ms.Stderr = nopWriterCloser{&prefixSuffixSaver{N: 1 << 25}}
	}

	go func() {
		defer ts.Cleanup(tb)
		err := ts.ms.Main(ctx, ts.sigs, ts.Run)
		if err != nil {
			if ts.Stderr == nil {
				err = fmt.Errorf("%w; stderr: %s", err, ts.ms.Stderr.(nopWriterCloser).Writer.(*prefixSuffixSaver).Bytes())
			}
		}
		ts.done <- err
	}()
}

func (ts *TestingState) Cleanup(tb testing.TB) {
	tb.Helper()

	if rc, ok := ts.Stdin.(io.ReadCloser); ok {
		err := rc.Close()
		if err != nil {
			tb.Errorf("failed to close stdin: %v", err)
		}
	}

	err, ok := ts.ExitError()
	if ok {
		// Already exited.
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = ts.Signal(ctx, os.Interrupt)
	if err != nil {
		tb.Errorf("failed to os.Interrupt testing xmain: %v", err)
	}
	err = ts.Wait(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		err = ts.Signal(ctx, os.Kill)
		if err != nil {
			tb.Errorf("failed to kill testing xmain: %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = ts.Wait(ctx)
	}
	assert.Success(tb, err)
}

func (ts *TestingState) Signal(ctx context.Context, sig os.Signal) (err error) {
	defer xdefer.Errorf(&err, "failed to signal testing xmain: %v", ts.ms.Name)

	err = ts.mu.Lock(ctx)
	if err != nil {
		return err
	}
	defer ts.mu.Unlock()

	if ts.doneErr != nil {
		return fmt.Errorf("testing xmain done: %w", *ts.doneErr)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ts.done:
		ts.doneErr = &err
		return err
	case ts.sigs <- sig:
		return nil
	}
}

func (ts *TestingState) Wait(ctx context.Context) (err error) {
	defer xdefer.Errorf(&err, "failed to wait testing xmain: %v", ts.ms.Name)

	err = ts.mu.Lock(ctx)
	if err != nil {
		return err
	}
	defer ts.mu.Unlock()

	if ts.doneErr != nil {
		if *ts.doneErr == nil {
			return nil
		}
		return fmt.Errorf("testing xmain done: %w", *ts.doneErr)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ts.done:
		ts.doneErr = &err
		return err
	}
}

func (ts *TestingState) ExitError() (error, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := ts.mu.Lock(ctx)
	if err != nil {
		return nil, false
	}
	defer ts.mu.Unlock()

	if ts.doneErr != nil {
		return *ts.doneErr, true
	}
	return nil, false
}

type nopWriterCloser struct {
	io.Writer
}

func (c nopWriterCloser) Close() error {
	return nil
}
