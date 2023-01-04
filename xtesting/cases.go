package xtesting

import (
	"context"
	"testing"
)

type Case struct {
	Name string
	Run  func(t *testing.T, ctx context.Context)
}

func RunCases(t *testing.T, ctx context.Context, tca []Case) {
	for _, tc := range tca {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			tc.Run(t, ctx)
		})
	}
}
