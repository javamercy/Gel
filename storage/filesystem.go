package storage

import (
	"Gel/core/constant"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type IFilesystemStorage interface {
	MakeDir(path string, permission os.FileMode) error
	WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error
	ReadFile(path string) ([]byte, error)
	Exists(path string) bool
}

type FilesystemStorage struct {
}

func NewFilesystemStorage() *FilesystemStorage {
	return &FilesystemStorage{}
}

func (filesystemStorage *FilesystemStorage) MakeDir(path string, permission os.FileMode) error {
	return os.MkdirAll(path, permission)
}

func (filesystemStorage *FilesystemStorage) WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error {
	if autoCreateDir {
		dir := filepath.Dir(path)
		if err := filesystemStorage.MakeDir(dir, constant.GelDirectoryPermission); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, permission)
}

func (filesystemStorage *FilesystemStorage) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}

func (filesystemStorage *FilesystemStorage) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
