package cmdlog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"oss.terrastruct.com/utils-go/assert"
	"oss.terrastruct.com/utils-go/xos"
	"oss.terrastruct.com/utils-go/xterm"
)

func init() {
	timeNow = func() time.Time {
		return time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC)
	}
}

func TestStdlog(t *testing.T) {
	tsw, ok := log.Default().Writer().(*tsWriter)
	if !ok {
		t.Fatalf("unexpected log.Default().Writer(): %T", log.Default().Writer())
	}
	b := &bytes.Buffer{}
	tsw.w = b

	log.Print("testing stdlog")
	exp := fmt.Sprintf("[01:01:01] %stesting stdlog\n",
		xterm.Prefix(xos.NewEnv(os.Environ()), os.Stderr, xterm.Blue, "stdlog"),
	)
	assert.String(t, exp, b.String())
}
