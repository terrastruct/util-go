package xterm

import (
	"fmt"
	"hash/crc64"
	"io"
	"math/rand"
	"os"

	"golang.org/x/term"

	"oss.terrastruct.com/util-go/xos"
)

// See https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Functions-using-CSI-_-ordered-by-the-final-character_s_
const (
	csi   = "\x1b["
	reset = csi + "0m"

	Bold = csi + "1m"

	Red     = csi + "31m"
	Green   = csi + "32m"
	Yellow  = csi + "33m"
	Blue    = csi + "34m"
	Magenta = csi + "35m"
	Cyan    = csi + "36m"

	BrightRed     = csi + "91m"
	BrightGreen   = csi + "92m"
	BrightYellow  = csi + "93m"
	BrightBlue    = csi + "94m"
	BrightMagenta = csi + "95m"
	BrightCyan    = csi + "96m"
)

var colors = [...]string{
	Red,
	Green,
	Yellow,
	Blue,
	Magenta,
	Cyan,

	BrightRed,
	BrightGreen,
	BrightYellow,
	BrightBlue,
	BrightMagenta,
	BrightCyan,
}

// isTTY checks whether the given writer is a *os.File TTY.
func isTTY(w io.Writer) bool {
	f, ok := w.(interface {
		Fd() uintptr
	})
	return ok && term.IsTerminal(int(f.Fd()))
}

func shouldColor(env *xos.Env, w io.Writer) bool {
	eb, err := env.Bool("COLOR")
	if eb != nil {
		return *eb
	}
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("xterm: %v", err))
	}
	if env.Getenv("TERM") == "dumb" {
		return false
	}
	return isTTY(w)
}

func Tput(env *xos.Env, w io.Writer, caps, s string) string {
	if caps == "" {
		return s
	}
	if !shouldColor(env, w) {
		return s
	}
	return caps + s + reset
}

func Prefix(env *xos.Env, w io.Writer, caps, s string) string {
	s = fmt.Sprintf("%s", s)
	return Tput(env, w, caps, s) + ":"
}

var crc64Table = crc64.MakeTable(crc64.ISO)

// CC meaning constant color. So constant color prefix.
func CCPrefix(env *xos.Env, w io.Writer, s string) string {
	sum := crc64.Checksum([]byte(s), crc64Table)
	rand := rand.New(rand.NewSource(int64(sum)))

	color := colors[rand.Intn(len(colors))]
	return Prefix(env, w, color, s)
}
