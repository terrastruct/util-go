package xcontext_test

import (
	"context"
	"testing"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/xcontext"
)

func TestWithoutCancel(t *testing.T) {
	t.Parallel()

	s := "meow"

	ctx := context.Background()
	ctx = stringWith(ctx, s)
	assert.Success(t, ctx.Err())

	t.Run("no_cancel", func(t *testing.T) {
		t.Parallel()

		ctx := xcontext.WithoutCancel(ctx)
		assert.Success(t, ctx.Err())

		s2 := stringFrom(ctx)
		assert.String(t, s, s2)
	})

	t.Run("cancel_before", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		assert.Error(t, ctx.Err())

		ctx = xcontext.WithoutCancel(ctx)
		assert.Success(t, ctx.Err())

		s2 := stringFrom(ctx)
		assert.String(t, s, s2)
	})

	t.Run("cancel_after", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(ctx)
		ctx = xcontext.WithoutCancel(ctx)
		cancel()
		assert.Success(t, ctx.Err())

		s2 := stringFrom(ctx)
		assert.String(t, s, s2)
	})
}

func TestWithoutValues(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const k = "Death is nature's way of saying `Howdy'."
	const exp = "Proposed Additions to the PDP-11 Instruction Set"
	const exp2 = "character density, n.:"

	t.Run("no_value", func(t *testing.T) {
		t.Parallel()

		v := ctx.Value(k)
		assert.JSON(t, nil, v)

		ctx := xcontext.WithoutValues(ctx)

		v = ctx.Value(k)
		assert.JSON(t, nil, v)
	})

	t.Run("with_value", func(t *testing.T) {
		t.Parallel()

		ctxv := context.WithValue(ctx, k, exp)
		ctx := xcontext.WithoutValues(ctxv)

		// ctxv contains k but ctx doesn't.
		v := ctxv.Value(k)
		assert.JSON(t, exp, v)

		v = ctx.Value(k)
		assert.JSON(t, nil, v)

		ctx = context.WithValue(ctx, k, exp2)
		v = ctx.Value(k)
		assert.JSON(t, exp2, v)
	})

	t.Run("cancel", func(t *testing.T) {
		t.Parallel()

		t.Run("before", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			cancel()

			ctx = xcontext.WithoutValues(ctx)
			assert.Error(t, ctx.Err())
		})

		t.Run("after", func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			ctx = xcontext.WithoutValues(ctx)
			assert.Success(t, ctx.Err())

			cancel()
			assert.Error(t, ctx.Err())
		})
	})
}

type stringKey struct{}

func stringFrom(ctx context.Context) string {
	return ctx.Value(stringKey{}).(string)
}

func stringWith(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, stringKey{}, s)
}
