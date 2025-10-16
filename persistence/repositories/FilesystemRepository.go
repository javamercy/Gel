package repositories

import (
	"errors"
	"io/fs"
	"os"
)

type FilesystemRepository struct {
}

func NewFilesystemRepository() *FilesystemRepository {
	return &FilesystemRepository{}
}

func (filesystemRepository *FilesystemRepository) MakeDirRange(paths []string) error {
	for _, path := range paths {
		if err := filesystemRepository.MakeDir(path); err != nil {
			return err
		}
	}
	return nil
}

func (filesystemRepository *FilesystemRepository) MakeDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (filesystemRepository *FilesystemRepository) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (filesystemRepository *FilesystemRepository) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !errors.Is(err, fs.ErrNotExist)
}

func (filesystemRepository *FilesystemRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
