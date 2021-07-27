package source

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/kyoh86/my-nerds/ioutil"
)

func CopyDirToLocal(source Source, localDir, src string) (retErr error) {
	return source.Walk(src, func(entry string, info fs.FileInfo) error {
		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel := strings.Split(strings.TrimLeft(strings.TrimPrefix(entry, src), "/"), "/")
		return CopyToLocal(source, filepath.Join(append([]string{localDir}, rel...)...), entry)
	})
}

func CopyToLocal(source Source, localFile, src string) (retErr error) {
	resp, err := source.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	return ioutil.CopyToFile(localFile, resp)
}
