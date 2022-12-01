package xtesting

import (
	"context"
	"testing"

	"oss.terrastruct.com/util-go/xos"
)

type Case struct {
	Name string
	Skip bool

	Run func(t *testing.T, ctx context.Context, env *xos.Env)
}

func RunCases(t *testing.T, cases []Case) {
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Skip {
				t.SkipNow()
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			env := xos.NewEnv(nil)

			tc.Run(t, ctx, env)
		})
	}
}
