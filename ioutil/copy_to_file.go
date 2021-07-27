package ioutil

import (
	"io"
	"os"
	"path/filepath"
)

func CopyToFile(filename string, reader io.Reader) (retErr error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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
