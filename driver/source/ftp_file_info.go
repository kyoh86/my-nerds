package source

import (
	"io/fs"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

type fileInfo struct {
	back *ftp.Walker

	once sync.Once
	stat *ftp.Entry
}

func (f *fileInfo) getStat() {
	f.once.Do(func() {
		f.stat = f.back.Stat()
	})
}

// Name is base name of the file
func (f *fileInfo) Name() string {
	f.getStat()
	return f.stat.Name
}

// Size in bytes for regular files; system-dependent for others
func (f *fileInfo) Size() int64 {
	f.getStat()
	return int64(f.stat.Size)
}

// Mode bits
func (f *fileInfo) Mode() fs.FileMode {
	f.getStat()
	switch f.stat.Type {
	case ftp.EntryTypeLink:
		return fs.ModeSymlink
	case ftp.EntryTypeFolder:
		return fs.ModeDir
	}
	return 0
}

// ModTime is modification time
func (f *fileInfo) ModTime() time.Time {
	f.getStat()
	return f.stat.Time
}

// IsDir is abbreviation for Mode().IsDir()
func (f *fileInfo) IsDir() bool {
	f.getStat()
	return f.stat.Type == ftp.EntryTypeFolder
}

// Sys is nil
func (f *fileInfo) Sys() interface{} {
	return f.back
}
