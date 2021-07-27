package source

import (
	"errors"
	"io"
	"io/fs"
)

var (
	ErrSkipDir     = errors.New("skip dir")
	ErrStopWalking = errors.New("stop walking")
)

type Source interface {
	Walk(root string, walkFn func(string, fs.FileInfo) error) error
	Rename(oldPath, newPath string) error
	Open(path string) (io.ReadCloser, error)
	Separator() string
}
