package repositories

import (
	"Gel/src/gel/core/constants"
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

func (filesystemRepository *FilesystemRepository) MakeDirRange(paths []string, permission os.FileMode) error {
	for _, path := range paths {
		if err := filesystemRepository.MakeDir(path, permission); err != nil {
			return err
		}
	}
	return nil
}

func (filesystemRepository *FilesystemRepository) MakeDir(path string, permission os.FileMode) error {
	return os.MkdirAll(path, permission)
}

func (filesystemRepository *FilesystemRepository) WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error {
	if autoCreateDir {
		dir := filepath.Dir(path)
		if err := filesystemRepository.MakeDir(dir, constants.DirPermission); err != nil {
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
	objectsDir, err := filesystemRepository.FindObjectsDir(startPath)
	if err != nil {
		return "", err
	}
	dir := hash[:2]
	file := hash[2:]
	objectPath := filepath.Join(objectsDir, dir, file)
	return objectPath, nil
}

func (filesystemRepository *FilesystemRepository) FindObjectsDir(startPath string) (string, error) {
	gelDir, err := filesystemRepository.FindGelDir(startPath)
	if err != nil {
		return "", err
	}
	objectsDir := filepath.Join(gelDir, constants.ObjectsDirName)
	return objectsDir, nil
}
