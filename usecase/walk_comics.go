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

func WalkComics(src source.Source, root string, walkFn func(model.Comic) error) error {
	done := map[string]struct{}{}
	return src.Walk(root, func(filename string, info fs.FileInfo) error {
		if info.IsDir() {
			return nil
		}

		comic := parseComicLikePath(filename)
		if comic == nil {
			return nil
		}
		if _, ok := done[comic.Path]; ok {
			return nil
		}
		done[comic.Path] = struct{}{}
		switch comic.Type {
		case model.ComicTypeRAR:
			name := path.Base(comic.Path)
			ext := path.Ext(name)
			name = strings.TrimSuffix(name, ext)
			i, err := model.ParseComicInfo(name)
			if err != nil {
				return nil
			}
			comic.Info = *i
		case model.ComicTypeDir:
			name := path.Base(comic.Path)
			i, err := model.ParseComicInfo(name)
			if err != nil {
				return nil
			}
			comic.Info = *i
		}
		if err := walkFn(*comic); err != nil {
			if errors.Is(err, ErrStopWalking) {
				return source.ErrStopWalking
			}
			return fmt.Errorf("found error when walking on %q: %w", filename, err)
		}
		return nil
	})
}

func parseComicLikePath(p string) *model.Comic {
	lower := strings.ToLower(p)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png") {
		return &model.Comic{Path: path.Dir(p), Type: model.ComicTypeDir}
	}
	if strings.HasSuffix(lower, ".rar") {
		return &model.Comic{Path: p, Type: model.ComicTypeRAR}
	}
	return nil
}
