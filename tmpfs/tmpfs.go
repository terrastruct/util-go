// Package tmpfs takes in an memory description of a filesystem
// and writes it to a temp directory so that it may be used as an io/fs.FS.
package tmpfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type FS struct {
	dir string
	fs.FS
}

func Make(m map[string]string) (*FS, error) {
	tempDir, err := os.MkdirTemp("", "tmpfs-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create root tmpfs dir: %w", err)
	}
	for p, s := range m {
		p = path.Join(tempDir, p)
		err = os.MkdirAll(path.Dir(p), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create tmpfs dir %q: %w", path.Dir(p), err)
		}
		err = os.WriteFile(p, []byte(s), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write tmpfs file %q: %w", p, err)
		}
	}
	return &FS{
		dir: tempDir,
		FS:  os.DirFS(tempDir),
	}, nil
}

func (fs *FS) Close() error {
	err := os.RemoveAll(fs.dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to close tmpfs.FS: %w", err)
	}
	return nil
}
