package ftp

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	host = "192.168.11.12:21"
	Root = "/sataraid1/nerd"
)

func Connect(ctx context.Context, user, pass string) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(host, ftp.DialWithTimeout(10*time.Second), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("dial to server: %w", err)
	}
	if err := conn.Login(user, pass); err != nil {
		return nil, fmt.Errorf("login on server: %w", err)
	}
	return conn, err
}

type fileInfo struct {
	back *ftp.Walker

	once sync.Once
	stat *ftp.Entry
}

func (f *fileInfo) getStat() {
	f.once.Do(func() {
		f.stat = f.back.Stat()
	})
}

// Name is base name of the file
func (f *fileInfo) Name() string {
	f.getStat()
	return f.stat.Name
}

// Size in bytes for regular files; system-dependent for others
func (f *fileInfo) Size() int64 {
	f.getStat()
	return int64(f.stat.Size)
}

// Mode bits
func (f *fileInfo) Mode() fs.FileMode {
	f.getStat()
	switch f.stat.Type {
	case ftp.EntryTypeLink:
		return fs.ModeSymlink
	case ftp.EntryTypeFolder:
		return fs.ModeDir
	}
	return 0
}

// ModTime is modification time
func (f *fileInfo) ModTime() time.Time {
	f.getStat()
	return f.stat.Time
}

// IsDir is abbreviation for Mode().IsDir()
func (f *fileInfo) IsDir() bool {
	f.getStat()
	return f.stat.Type == ftp.EntryTypeFolder
}

// Sys is nil
func (f *fileInfo) Sys() interface{} {
	return nil
}

var (
	ErrSkipDir     = errors.New("skip dir")
	ErrStopWalking = errors.New("stop walking")
)

// Walk on the path
func Walk(conn *ftp.ServerConn, root string, walkFn func(string, fs.FileInfo) error) error {
	walker := conn.Walk(root)
	for walker.Next() {
		p := walker.Path()
		// Ignore apple double files
		if path.Base(p) == ".AppleDouble" || path.Base(path.Dir(p)) == ".AppleDouble" {
			walker.SkipDir()
			continue
		}
		if err := walkFn(p, &fileInfo{back: walker}); err != nil {
			if errors.Is(err, ErrSkipDir) {
				walker.SkipDir()
				continue
			}
			if errors.Is(err, ErrStopWalking) {
				return nil
			}
			return err
		}
	}
	return walker.Err()
}
