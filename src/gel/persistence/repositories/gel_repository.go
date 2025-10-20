package repositories

import (
	"Gel/src/gel/core/constants"
	"errors"
	"path/filepath"
)

type IGelRepository interface {
	FindGelDir(startPath string) (string, error)
	FindObjectsDir(startPath string) (string, error)
	FindIndexFilePath(startPath string) (string, error)
	FindObjectPath(hash string, startPath string) (string, error)
	WriteObject(hash string, data []byte, startPath string) error
	ReadObject(hash string, startPath string) ([]byte, error)
	ObjectExists(hash string, startPath string) bool
}

type GelRepository struct {
	filesystemRepo IFilesystemRepository
}

func NewGelRepository(filesystemRepo IFilesystemRepository) *GelRepository {
	return &GelRepository{
		filesystemRepo: filesystemRepo,
	}
}

func (gelRepository *GelRepository) FindGelDir(startPath string) (string, error) {
	currentPath := startPath
	for {
		gelPath := filepath.Join(currentPath, constants.RepositoryDirName)
		if gelRepository.filesystemRepo.Exists(gelPath) {
			return gelPath, nil
		}
		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			break
		}
		currentPath = parent
	}
	return "", errors.New("not a gel repository (or any of the parent directories)")
}

func (gelRepository *GelRepository) FindObjectsDir(startPath string) (string, error) {
	gelDir, err := gelRepository.FindGelDir(startPath)
	if err != nil {
		return "", err
	}
	objectsDir := filepath.Join(gelDir, constants.ObjectsDirName)
	return objectsDir, nil
}

func (gelRepository *GelRepository) FindIndexFilePath(startPath string) (string, error) {
	gelDir, err := gelRepository.FindGelDir(startPath)
	if err != nil {
		return "", err
	}
	indexFilePath := filepath.Join(gelDir, constants.IndexFileName)
	return indexFilePath, nil
}

func (gelRepository *GelRepository) FindObjectPath(hash string, startPath string) (string, error) {
	objectsDir, err := gelRepository.FindObjectsDir(startPath)
	if err != nil {
		return "", err
	}
	dir := hash[:2]
	file := hash[2:]
	objectPath := filepath.Join(objectsDir, dir, file)
	return objectPath, nil
}

func (gelRepository *GelRepository) WriteObject(hash string, data []byte, startPath string) error {
	objectPath, err := gelRepository.FindObjectPath(hash, startPath)
	if err != nil {
		return err
	}
	return gelRepository.filesystemRepo.WriteFile(objectPath, data, true, constants.FilePermission)
}

func (gelRepository *GelRepository) ReadObject(hash string, startPath string) ([]byte, error) {
	objectPath, err := gelRepository.FindObjectPath(hash, startPath)
	if err != nil {
		return nil, err
	}
	return gelRepository.filesystemRepo.ReadFile(objectPath)
}

func (gelRepository *GelRepository) ObjectExists(hash string, startPath string) bool {
	objectPath, err := gelRepository.FindObjectPath(hash, startPath)
	if err != nil {
		return false
	}
	return gelRepository.filesystemRepo.Exists(objectPath)
}
