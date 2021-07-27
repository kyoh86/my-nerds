package archive

import (
	"io"
	"path/filepath"

	"github.com/gen2brain/go-unarr"
	"github.com/kyoh86/my-nerds/ioutil"
)

func ExtractRAR(pathTo string, reader io.Reader) (retErr error) {
	arch, err := unarr.NewArchiveFromReader(reader)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := arch.Close(); closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()
	entries, err := arch.List()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := arch.EntryFor(entry); err != nil {
			return err
		}
		if err := ioutil.CopyToFile(filepath.Join(pathTo, arch.Name()), arch); err != nil {
			return err
		}
	}
	return nil
}
