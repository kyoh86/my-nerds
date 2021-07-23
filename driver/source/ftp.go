package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPServer struct {
	conn *ftp.ServerConn
}

func OpenFTPServer(ctx context.Context, host, user, pass string) (*FTPServer, error) {
	conn, err := ftp.Dial(host, ftp.DialWithTimeout(10*time.Second), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("dial to server: %w", err)
	}
	if err := conn.Login(user, pass); err != nil {
		return nil, fmt.Errorf("login on server: %w", err)
	}
	return &FTPServer{conn}, err
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
func (s *FTPServer) Walk(root string, walkFn func(string, fs.FileInfo) error) error {
	walker := s.conn.Walk(root)
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

func (s *FTPServer) Rename(oldPath, newPath string) error {
	return s.conn.Rename(oldPath, newPath)
}

func (s *FTPServer) List(p string, walkFn func(string) error) error {
	entries, err := s.conn.List(p)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Type != ftp.EntryTypeFile {
			continue
		}
		if err := walkFn(path.Join(p, entry.Name)); err != nil {
			return err
		}
	}
	return nil
}

func (s *FTPServer) Download(pathFrom, pathTo string) (retErr error) {
	resp, err := s.Open(pathFrom)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return save(pathTo, resp)
}

func (s *FTPServer) Open(path string) (_ io.ReadCloser, retErr error) {
	resp, err := s.conn.Retr(path)
	if err != nil {
		return nil, err
	}
	if err := resp.SetDeadline(time.Now().Add(10 * time.Minute)); err != nil {
		resp.Close()
		return nil, err
	}
	return resp, nil
}

func save(pathTo string, reader io.Reader) (retErr error) {
	file, err := os.OpenFile(pathTo, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}
	return nil
}
