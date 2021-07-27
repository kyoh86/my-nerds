package usecase

import (
	"github.com/kyoh86/my-nerds/driver/archive"
	"github.com/kyoh86/my-nerds/driver/source"
)

func DownloadComicDir(src source.Source, pathFrom, pathTo string) (retErr error) {
	return source.CopyDirToLocal(src, pathTo, pathFrom)
}

func ExtractComic(src source.Source, arch archive.Extractor, pathFrom, pathTo string) (retErr error) {
	rar, err := src.Open(pathFrom)
	if err != nil {
		return err
	}
	defer func() {
		if err := rar.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return arch(pathTo, rar)
}
