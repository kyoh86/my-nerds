package source

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Local struct {
	root string
}

func OpenLocal() (*Local, error) {
	return new(Local), nil
}

func (l *Local) Walk(dir string, walkFn func(string, fs.FileInfo) error) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return walkFn(path, info)
	})
}

func (l *Local) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

func (l *Local) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *Local) Separator() string {
	return string([]rune{filepath.Separator})
}
