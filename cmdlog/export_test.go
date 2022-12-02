package cmdlog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/xos"
	"oss.terrastruct.com/util-go/xterm"
)

func init() {
	timeNow = func() time.Time {
		return time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC)
	}
}

func TestStdlog(t *testing.T) {
	pw, ok := log.Default().Writer().(prefixWriter)
	if !ok {
		t.Fatalf("unexpected log.Default().Writer(): %T", log.Default().Writer())
	}
	tsw, ok := pw.w.(*tsWriter)
	if !ok {
		t.Fatalf("unexpected pw.w: %T", pw.w)
	}
	b := &bytes.Buffer{}
	ow := tsw.w
	tsw.w = b
	defer func() {
		tsw.w = ow
	}()

	log.Print("testing stdlog")
	exp := fmt.Sprintf("[01:01:01] %s testing stdlog\n",
		xterm.Prefix(xos.NewEnv(os.Environ()), os.Stderr, xterm.Blue, "stdlog"),
	)
	assert.String(t, exp, b.String())
}
