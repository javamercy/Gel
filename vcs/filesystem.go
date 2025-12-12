package vcs

import (
	"Gel/storage"
	"os"
)

type FilesystemService struct {
	filesystemStorage storage.IFilesystemStorage
}

func NewFilesystemService(filesystemStorage storage.IFilesystemStorage) *FilesystemService {
	return &FilesystemService{
		filesystemStorage: filesystemStorage,
	}
}

func (filesystemService *FilesystemService) ReadFile(path string) ([]byte, error) {
	return filesystemService.filesystemStorage.ReadFile(path)
}

func (filesystemService *FilesystemService) WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error {
	return filesystemService.filesystemStorage.WriteFile(path, data, autoCreateDir, permission)
}

func (filesystemService *FilesystemService) Exists(path string) bool {
	return filesystemService.filesystemStorage.Exists(path)
}

func (filesystemService *FilesystemService) MakeDirectory(path string, permission os.FileMode) error {
	return filesystemService.filesystemStorage.MakeDir(path, permission)
}
