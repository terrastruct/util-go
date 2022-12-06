package xexec

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// findExecutable is from package os/exec
func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	m := d.Mode()
	if !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return fs.ErrPermission
}

// SearchPath searches for all executables that have prefix in their names in
// the directories named by the PATH environment variable.
func SearchPath(prefix string) ([]string, error) {
	var matches []string
	envPath := os.Getenv("PATH")
	dirSet := make(map[string]struct{})
	for _, dir := range filepath.SplitList(envPath) {
		if dir == "" {
			// From exec package:
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		if _, ok := dirSet[dir]; ok {
			continue
		}
		dirSet[dir] = struct{}{}
		files, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		for _, f := range files {
			if strings.HasPrefix(f.Name(), prefix) {
				match := filepath.Join(dir, f.Name())
				// Unideal but I don't want to maintain two separate implementations of this
				// function like os/exec.
				if runtime.GOOS == "windows" {
					matches = append(matches, match)
					continue
				}
				err = findExecutable(match)
				if err == nil {
					matches = append(matches, match)
				}
			}
		}

	}
	return matches, nil
}
