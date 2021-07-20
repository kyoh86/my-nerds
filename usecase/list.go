package usecase

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/kyoh86/my-nerds/driver/ftp"
)

type ComicType int

const (
	ComicTypeUnknown ComicType = iota
	ComicTypeRAR
	ComicTypeDir
)

var (
	ErrStopWalking = errors.New("stop walking")
)

func WalkComics(ctx context.Context, user, pass string, walkFn func(string, ComicType) error) error {
	conn, err := ftp.Connect(ctx, user, pass)
	if err != nil {
		return err
	}
	done := map[string]struct{}{}
	return ftp.Walk(conn, ftp.Root, func(path string, info fs.FileInfo) error {
		if info.IsDir() {
			return nil
		}

		p, t, ok := parseComicPath(path)
		if !ok {
			return nil
		}
		if _, ok := done[p]; ok {
			return nil
		}
		if err := walkFn(p, t); err != nil {
			if errors.Is(err, ErrStopWalking) {
				return ftp.ErrStopWalking
			}
			return fmt.Errorf("found comic %q: %w", p, err)
		}
		return nil
	})
}

func parseComicPath(p string) (string, ComicType, bool) {
	lower := strings.ToLower(p)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png") {
		return path.Dir(p), ComicTypeDir, true
	}
	if strings.HasSuffix(lower, ".rar") {
		return p, ComicTypeRAR, true
	}
	return p, 0, false
}
