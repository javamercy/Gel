package repositories

import (
	"Gel/core/constants"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

func (filesystemRepository *FilesystemRepository) FindGelDir(startPath string) (string, error) {
	for {
		gelPath := filepath.Join(startPath, constants.RepositoryDirName)
		if filesystemRepository.Exists(gelPath) {
			return gelPath, nil
		}
		parent := filepath.Dir(startPath)
		if parent == startPath {
			break
		}
		startPath = parent
	}
	return "", errors.New("not a gel repository (or any of the parent directories)")
}

func (filesystemRepository *FilesystemRepository) FindObjectPath(hash string, startPath string) (string, error) {
	gelDir, err := filesystemRepository.FindGelDir(startPath)
	if err != nil {
		return "", err
	}
	objectDir := filepath.Join(gelDir, constants.ObjectsDirName, hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])
	if !filesystemRepository.Exists(objectPath) {
		return "", errors.New("object not found")
	}
	return objectPath, nil
}
