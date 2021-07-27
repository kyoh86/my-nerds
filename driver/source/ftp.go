package source

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
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

// Walk on the path
func (s *FTPServer) Walk(dir string, walkFn func(string, fs.FileInfo) error) error {
	walker := s.conn.Walk(dir)
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

func (s *FTPServer) Separator() string { return "/" }

var _ Source = (*FTPServer)(nil)
