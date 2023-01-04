package xtesting

import (
	"context"
	"testing"

	"oss.terrastruct.com/util-go/xos"
)

type Case struct {
	Name string

	Run func(t *testing.T, ctx context.Context, env *xos.Env)
}

func RunCases(t *testing.T, cases []Case) {
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			env := xos.NewEnv(nil)

			tc.Run(t, ctx, env)
		})
	}
}
