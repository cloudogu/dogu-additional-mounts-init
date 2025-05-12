package copy

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Filesystem interface {
	Lstat(path string) (os.FileInfo, error)
	EvalSymlinks(path string) (string, error)
	Stat(name string) (os.FileInfo, error)
	Open(name string) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	CloseFile(file *os.File) error
	SyncFile(file *os.File) error
	SameFile(fi1, fi2 os.FileInfo) bool
	WalkDir(root string, fn fs.WalkDirFunc) error
}

type fileSystem struct{}

func (f fileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

func (f fileSystem) SameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

func (f fileSystem) SyncFile(file *os.File) error {
	return file.Sync()
}

func (f fileSystem) CloseFile(file *os.File) error {
	return file.Close()
}

func (f fileSystem) Lstat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

func (f fileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (f fileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f fileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (f fileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f fileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (f fileSystem) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
