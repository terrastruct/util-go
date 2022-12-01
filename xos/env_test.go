package xos_test

import (
	"testing"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/xos"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	e := xos.NewEnv(nil)
	assert.String(t, "", e.Getenv("NONE"))

	e.Setenv("DEBUG", "MEOW")
	e.Setenv("DEBUG2", "TWO")
	e.Setenv("DEBUG3", "THREE")
	e.Setenv("DEBUG4", "FOUR")

	assert.StringJSON(t, `[
  "DEBUG=MEOW",
  "DEBUG2=TWO",
  "DEBUG3=THREE",
  "DEBUG4=FOUR"
]`, e.Environ())

	e.Setenv("DEBUG", "D2")
	assert.StringJSON(t, `[
  "DEBUG=D2",
  "DEBUG2=TWO",
  "DEBUG3=THREE",
  "DEBUG4=FOUR"
]`, e.Environ())
	assert.String(t, "D2", e.Getenv("DEBUG"))
}

func TestBool(t *testing.T) {
	t.Parallel()

	env := xos.NewEnv(nil)

	eb, err := env.Bool("MY_BOOL")
	assert.Success(t, err)
	assert.StringJSON(t, `null`, eb)

	env.Setenv("MY_BOOL", "0")
	eb, err = env.Bool("MY_BOOL")
	assert.Success(t, err)
	assert.StringJSON(t, `false`, eb)

	env.Setenv("MY_BOOL", "1")
	eb, err = env.Bool("MY_BOOL")
	assert.Success(t, err)
	assert.StringJSON(t, `true`, eb)

	env.Setenv("MY_BOOL", "false")
	eb, err = env.Bool("MY_BOOL")
	assert.Success(t, err)
	assert.StringJSON(t, `false`, eb)

	env.Setenv("MY_BOOL", "true")
	eb, err = env.Bool("MY_BOOL")
	assert.Success(t, err)
	assert.StringJSON(t, `true`, eb)

	env.Setenv("MY_BOOL", "TRUE")
	eb, err = env.Bool("MY_BOOL")
	assert.Error(t, err)
	assert.ErrorContains(t, err, `$MY_BOOL must be 0, 1, false or true but got "TRUE"`)
}
