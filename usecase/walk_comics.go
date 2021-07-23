package usecase

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/kyoh86/my-nerds/driver/source"
	"github.com/kyoh86/my-nerds/model"
)

var (
	ErrStopWalking = errors.New("stop walking")
)

func WalkComics(server *source.FTPServer, root string, walkFn func(string, model.ComicType) error) error {
	done := map[string]struct{}{}
	return server.Walk(root, func(path string, info fs.FileInfo) error {
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
		done[p] = struct{}{}
		if err := walkFn(p, t); err != nil {
			if errors.Is(err, ErrStopWalking) {
				return source.ErrStopWalking
			}
			return fmt.Errorf("found comic %q: %w", p, err)
		}
		return nil
	})
}

func parseComicPath(p string) (string, model.ComicType, bool) {
	lower := strings.ToLower(p)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png") {
		return path.Dir(p), model.ComicTypeDir, true
	}
	if strings.HasSuffix(lower, ".rar") {
		return p, model.ComicTypeRAR, true
	}
	return p, 0, false
}
