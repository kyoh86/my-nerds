package usecase

import "github.com/kyoh86/my-nerds/driver/source"

func ExtractComicRAR(server *source.FTPServer, pathFrom, pathTo string) (retErr error) {
	rar, err := server.Open(pathFrom)
	if err != nil {
		return err
	}
	defer func() {
		if err := rar.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return source.ExtractRAR(rar, pathTo)
}
