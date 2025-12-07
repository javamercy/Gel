package repositories

import (
	"Gel/core/constant"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type IFilesystemRepository interface {
	MakeDir(path string, permission os.FileMode) error
	WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error
	ReadFile(path string) ([]byte, error)
	Exists(path string) bool
}

type FilesystemRepository struct {
}

func NewFilesystemRepository() *FilesystemRepository {
	return &FilesystemRepository{}
}

func (filesystemRepository *FilesystemRepository) MakeDir(path string, permission os.FileMode) error {
	return os.MkdirAll(path, permission)
}

func (filesystemRepository *FilesystemRepository) WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error {
	if autoCreateDir {
		dir := filepath.Dir(path)
		if err := filesystemRepository.MakeDir(dir, constant.GelDirPermission); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, permission)
}

func (filesystemRepository *FilesystemRepository) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}

func (filesystemRepository *FilesystemRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
