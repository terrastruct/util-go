package xos

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

type Env struct {
	environMu sync.RWMutex
	environ   []string
}

func NewEnv(environ []string) *Env {
	return &Env{
		environ: environ,
	}
}

func (e *Env) Environ() []string {
	e.environMu.RLock()
	defer e.environMu.RUnlock()

	environ2 := make([]string, 0, len(e.environ))
	return append(environ2, e.environ...)
}

func (e *Env) Setenv(name, value string) {
	e.environMu.Lock()
	defer e.environMu.Unlock()

	l := fmt.Sprintf("%s=%s", name, value)

	for i, l2 := range e.environ {
		j := strings.Index(l2, "=")
		if j == -1 {
			continue
		}
		name2 := l2[:j]
		if name != name2 {
			continue
		}
		e.environ[i] = l
		return
	}
	e.environ = append(e.environ, l)
}

func (e *Env) Getenv(name string) string {
	e.environMu.RLock()
	defer e.environMu.RUnlock()

	for _, l := range e.environ {
		i := strings.Index(l, "=")
		if i == -1 {
			continue
		}
		name2 := l[:i]
		if name == name2 {
			return l[i+1:]
		}
	}
	return ""
}

func (e *Env) Bool(name string) (*bool, error) {
	ev := e.Getenv(name)
	if ev == "" {
		return nil, nil
	}
	eb := new(bool)
	if ev == "0" || ev == "false" {
		return eb, nil
	}
	if ev == "1" || ev == "true" {
		*eb = true
		return eb, nil
	}
	return nil, fmt.Errorf("$%s must be 0, 1, false or true but got %q", name, ev)
}

func (e *Env) HumanPath(fp string) string {
	if strings.HasPrefix(fp, e.Getenv("HOME")) {
		return filepath.Join("~", strings.TrimPrefix(fp, e.Getenv("HOME")))
	}
	return fp
}
