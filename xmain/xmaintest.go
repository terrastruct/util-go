package xmain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/cmdlog"
	"oss.terrastruct.com/util-go/xdefer"
	"oss.terrastruct.com/util-go/xos"
)

type TestState struct {
	Run  func(context.Context, *State) error
	Env  *xos.Env
	Args []string
	PWD  string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	ms      *State
	sigs    chan os.Signal
	done    chan struct{}
	doneErr *error
}

func (ts *TestState) StdinPipe() (pw io.WriteCloser) {
	ts.Stdin, pw = io.Pipe()
	return pw
}

func (ts *TestState) StdoutPipe() (pr io.Reader) {
	pr, ts.Stdout = io.Pipe()
	return pr
}

func (ts *TestState) StderrPipe() (pr io.Reader) {
	pr, ts.Stderr = io.Pipe()
	return pr
}

func (ts *TestState) Start(tb testing.TB, ctx context.Context) {
	tb.Helper()

	if ts.done != nil {
		tb.Fatal("xmain.TestState.Start cannot be called twice")
	}

	if ts.Env == nil {
		ts.Env = xos.NewEnv(nil)
	}
	var tempDirCleanup func()
	if ts.PWD == "" {
		ts.PWD, tempDirCleanup = assert.TempDir(tb)
	}

	ts.sigs = make(chan os.Signal, 1)
	ts.done = make(chan struct{})

	name := ""
	args := []string(nil)
	if len(ts.Args) > 0 {
		name = ts.Args[0]
		args = ts.Args[1:]
	}
	log := cmdlog.NewTB(ts.Env, tb)
	ts.ms = &State{
		Name: name,

		Log:  log,
		Env:  ts.Env,
		Opts: NewOpts(ts.Env, log, args),
		PWD:  ts.PWD,
	}

	if ts.Stdin == nil {
		ts.ms.Stdin = io.LimitReader(nil, 0)
	} else if rc, ok := ts.Stdin.(io.ReadCloser); ok {
		ts.ms.Stdin = rc
	} else {
		var pw io.Writer
		ts.ms.Stdin, pw = io.Pipe()
		go io.Copy(pw, ts.Stdin)
	}

	var pipeWG sync.WaitGroup
	if ts.Stdout == nil {
		ts.ms.Stdout = nopWriterCloser{io.Discard}
	} else if wc, ok := ts.Stdout.(io.WriteCloser); ok {
		ts.ms.Stdout = wc
	} else {
		var pr io.Reader
		pr, ts.ms.Stdout = io.Pipe()
		pipeWG.Add(1)
		go func() {
			defer pipeWG.Done()
			io.Copy(ts.Stdout, pr)
		}()
	}
	if ts.Stderr == nil {
		ts.ms.Stderr = nopWriterCloser{&prefixSuffixSaver{N: 1 << 25}}
	} else if wc, ok := ts.Stderr.(io.WriteCloser); ok {
		ts.ms.Stderr = wc
	} else {
		var pr io.Reader
		pr, ts.ms.Stderr = io.Pipe()
		pipeWG.Add(1)
		go func() {
			defer pipeWG.Done()
			io.Copy(ts.Stderr, pr)
		}()
	}

	go func() {
		var err error
		defer func() {
			ts.ms.Stdout.Close()
			ts.ms.Stderr.Close()
			pipeWG.Wait()
			if tempDirCleanup != nil {
				tempDirCleanup()
			}
			ts.doneErr = &err
			close(ts.done)
			ts.cleanup(tb)
		}()
		err = ts.ms.Main(ctx, ts.sigs, ts.Run)
		if err != nil {
			if ts.Stderr == nil {
				stderr := ts.ms.Stderr.(nopWriterCloser).Writer.(*prefixSuffixSaver).Bytes()
				if len(stderr) > 0 {
					err = fmt.Errorf("%w; stderr: %s", err, stderr)
				}
			}
		}
	}()
}

func (ts *TestState) cleanup(tb testing.TB) {
	tb.Helper()
	if rc, ok := ts.ms.Stdin.(io.ReadCloser); ok {
		err := rc.Close()
		if err != nil {
			tb.Errorf("failed to close xmain test stdin: %v", err)
		}
	}
}

func (ts *TestState) Cleanup(tb testing.TB) {
	tb.Helper()

	ts.cleanup(tb)

	select {
	case <-ts.done:
		// Already exited.
		return
	default:
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := ts.Signal(ctx, os.Interrupt)
	if err != nil {
		tb.Errorf("failed to os.Interrupt xmain test: %v", err)
	}
	err = ts.Wait(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		err = ts.Signal(ctx, os.Kill)
		if err != nil {
			tb.Errorf("failed to kill xmain test: %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = ts.Wait(ctx)
	}
	assert.Success(tb, err)
}

func (ts *TestState) Signal(ctx context.Context, sig os.Signal) (err error) {
	defer xdefer.Errorf(&err, "failed to signal xmain test: %v", ts.ms.Name)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ts.done:
		return fmt.Errorf("xmain test exited: %w", *ts.doneErr)
	case ts.sigs <- sig:
		return nil
	}
}

func (ts *TestState) Wait(ctx context.Context) (err error) {
	defer xdefer.Errorf(&err, "failed to wait xmain test: %v", ts.ms.Name)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ts.done:
		return *ts.doneErr
	}
}

type nopWriterCloser struct {
	io.Writer
}

func (c nopWriterCloser) Close() error {
	return nil
}
