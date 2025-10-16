package repositories

import "os"

type FilesystemRepository struct {
}

func NewFilesystemRepository() *FilesystemRepository {
	return &FilesystemRepository{}
}

func (filesystemRepository *FilesystemRepository) MakeDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func (filesystemRepository *FilesystemRepository) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (filesystemRepository *FilesystemRepository) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (filesystemRepository *FilesystemRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
