package source

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-unarr"
)

func ExtractRAR(reader io.Reader, pathTo string) (retErr error) {
	arch, err := unarr.NewArchiveFromReader(reader)
	if err != nil {
		return err
	}
	defer func() {
		if err := arch.Close(); err != nil && retErr == nil {
			retErr = err
		}
	}()
	entries, err := arch.List()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(pathTo, 0755); err != nil {
		return err
	}
	for _, entry := range entries {
		if err := arch.EntryFor(entry); err != nil {
			return err
		}
		if err := save(filepath.Join(pathTo, arch.Name()), arch); err != nil {
			return err
		}
	}
	return nil
}
