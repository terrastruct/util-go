package xmain

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"
	"context"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/cmdlog"
	"oss.terrastruct.com/util-go/xos"
)

type TestingState struct {
	State  *State
	Stdin  io.Writer
	Stdout io.Reader
	Stderr io.Reader

	sigs chan os.Signal
	done chan error
}

func (ts *TestingState) Signal(ctx context.Context, sig os.Signal) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ts.sigs <- sig:
		return nil
	}
}

func (ts *TestingState) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ts.done:
		return err
	}
}

func Testing(tb testing.TB, ctx context.Context, env *xos.Env, run func(context.Context, *State) error, name string, args ...string) (ts *TestingState, cleanup func()) {
	stdinr, stdinw, err := os.Pipe()
	assert.Success(tb, err)
	stdoutr, stdoutw, err := os.Pipe()
	assert.Success(tb, err)
	stderrr, stderrw, err := os.Pipe()
	assert.Success(tb, err)

	ms := &State{
		Name: name,

		Stdin:  stdinr,
		Stdout: stdoutw,
		Stderr: stderrw,

		Env: env,
		Log: cmdlog.NewTB(env, tb),
	}
	ms.Opts = NewOpts(ms.Env, ms.Log, args)

	ts = &TestingState{
		State: ms,

		Stdin:  stdinw,
		Stdout: stdoutr,
		Stderr: stderrr,

		sigs: make(chan os.Signal, 1),
		done: make(chan error, 1),
	}

	cleanup = func() {
		stdinr.Close()
		stdinw.Close()
		stdoutr.Close()
		stdoutw.Close()
		stderrr.Close()
		stderrw.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err = ts.Signal(ctx, os.Interrupt)
		if err != nil {
			tb.Errorf("failed to os.Interrupt testing xmain: %v", err)
		}
		err := ts.Wait(ctx)
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

	go func() {
		ts.done <- ms.Main(ctx, ts.sigs, run)
	}()

	return ts, cleanup
}
