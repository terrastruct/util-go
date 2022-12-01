// Package xcontext implements indispensable context helpers.
package xcontext

import (
	"context"
	"time"
)

// WithoutCancel returns a context derived from ctx that may
// never be cancelled.
func WithoutCancel(ctx context.Context) context.Context {
	return withoutCancel{ctx: ctx}
}

type withoutCancel struct {
	ctx context.Context
}

func (c withoutCancel) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c withoutCancel) Done() <-chan struct{} {
	return nil
}

func (c withoutCancel) Err() error {
	return nil
}

func (c withoutCancel) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// WithoutValues creates a new context derived from ctx that does not inherit its values
// but does pass on cancellation.
func WithoutValues(ctx context.Context) context.Context {
	return withoutValues{ctx: ctx}
}

type withoutValues struct {
	ctx context.Context
}

func (c withoutValues) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c withoutValues) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c withoutValues) Err() error {
	return c.ctx.Err()
}

func (c withoutValues) Value(key interface{}) interface{} {
	return nil
}
