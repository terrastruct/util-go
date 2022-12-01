// Package cmdlog implements color leveled logging for command line tools.
package cmdlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"oss.terrastruct.com/utils-go/xos"
	"oss.terrastruct.com/utils-go/xterm"
)

var timeNow = time.Now

const defaultTSFormat = "15:04:05"

func init() {
	l := New(xos.NewEnv(os.Environ()), os.Stderr)
	l.SetTS(true)
	l = l.WithPrefix(xterm.Blue, "stdlog")

	log.SetOutput(l.NoLevel.Writer())
	log.SetPrefix(l.NoLevel.Prefix())
	log.SetFlags(l.NoLevel.Flags())
}

type Logger struct {
	env *xos.Env
	w   io.Writer
	tsw *tsWriter
	dw  *debugWriter

	NoLevel *log.Logger
	Debug   *log.Logger
	Success *log.Logger
	Info    *log.Logger
	Warn    *log.Logger
	Error   *log.Logger
}

func (l *Logger) GetTS() bool {
	l.tsw.mu.Lock()
	defer l.tsw.mu.Unlock()
	return l.tsw.enabled
}

func (l *Logger) GetTSFormat() string {
	l.tsw.mu.Lock()
	defer l.tsw.mu.Unlock()
	return l.tsw.tsfmt
}

func (l *Logger) GetDebug() bool {
	return l.dw.debug()
}

func (l *Logger) SetTS(enabled bool) {
	l.tsw.mu.Lock()
	l.tsw.enabled = enabled
	l.tsw.mu.Unlock()
}

func (l *Logger) SetTSFormat(tsfmt string) {
	l.tsw.mu.Lock()
	l.tsw.tsfmt = tsfmt
	l.tsw.mu.Unlock()
}

func (l *Logger) SetDebug(enabled bool) {
	vi := int64(0)
	if enabled {
		vi = 1
	}
	atomic.StoreInt64(&l.dw.flag, vi)
}

func New(env *xos.Env, w io.Writer) *Logger {
	tsw := &tsWriter{w: w, tsfmt: defaultTSFormat}
	dw := &debugWriter{w: tsw, env: env}
	l := &Logger{
		env: env,
		w:   w,
		dw:  dw,
		tsw: tsw,
	}
	l.init("")
	return l
}

func (l *Logger) init(prefix string) {
	l.NoLevel = log.New(l.tsw, prefix, 0)
	l.Debug = log.New(l.dw, prefix+xterm.Prefix(l.env, l.w, "", "debug"), 0)
	l.Success = log.New(l.tsw, prefix+xterm.Prefix(l.env, l.w, xterm.Green, "success"), 0)
	l.Info = log.New(l.tsw, prefix+xterm.Prefix(l.env, l.w, xterm.Blue, "info"), 0)
	l.Warn = log.New(l.tsw, prefix+xterm.Prefix(l.env, l.w, xterm.Yellow, "warn"), 0)
	l.Error = log.New(l.tsw, prefix+xterm.Prefix(l.env, l.w, xterm.Red, "err"), 0)
}

type debugWriter struct {
	w    io.Writer
	flag int64
	env  *xos.Env
}

func (dw *debugWriter) debug() bool {
	if atomic.LoadInt64(&dw.flag) == 0 {
		debug, err := dw.env.Bool("DEBUG")
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("cmdlog: %v", err))
			return false
		}
		if debug == nil || !*debug {
			return false
		}
	}
	return true
}

func (dw *debugWriter) Write(p []byte) (int, error) {
	if !dw.debug() {
		return len(p), nil
	}
	return dw.w.Write(p)
}

type tsWriter struct {
	w io.Writer

	mu      sync.Mutex
	tsfmt   string
	enabled bool
}

func (tsw *tsWriter) Write(p []byte) (int, error) {
	tsw.mu.Lock()
	enabled := tsw.enabled
	tsfmt := tsw.tsfmt
	tsw.mu.Unlock()

	if !enabled {
		return tsw.w.Write(p)
	}

	ts := timeNow().Format(tsfmt)
	prefix := []byte("[" + ts + "] ")
	p = append(prefix, p...)
	n, err := tsw.w.Write(p)
	if err != nil {
		n -= len(prefix)
		if n < 0 {
			n = 0
		}
		return n, err
	}
	return len(p), nil
}

func NewTB(env *xos.Env, tb testing.TB) *Logger {
	return New(env, tbWriter{tb})
}

type tbWriter struct {
	tb testing.TB
}

func (tbw tbWriter) Write(p []byte) (int, error) {
	tbw.tb.Logf("%s", p)
	return len(p), nil
}

// Allows detection as a terminal.
func (tbWriter) Fd() uintptr {
	return os.Stderr.Fd()
}

func (l *Logger) WithCCPrefix(s string) *Logger {
	l2 := new(Logger)
	*l2 = *l

	prefix := l.NoLevel.Prefix() + xterm.CCPrefix(l.env, l.w, s)
	l2.init(prefix)
	return l2
}

func (l *Logger) WithPrefix(caps, s string) *Logger {
	l2 := new(Logger)
	*l2 = *l

	prefix := l.NoLevel.Prefix() + xterm.Prefix(l.env, l.w, caps, s)
	l2.init(prefix)
	return l2
}
