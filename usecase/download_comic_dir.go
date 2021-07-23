package usecase

import (
	"path"
	"path/filepath"

	"github.com/kyoh86/my-nerds/driver/source"
)

func DownloadComicDir(server *source.FTPServer, pathFrom, pathTo string) (retErr error) {
	return server.List(pathFrom, func(entry string) error {
		return server.Download(entry, filepath.Join(pathTo, path.Base(entry)))
	})
}
