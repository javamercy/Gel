package services

import (
	"Gel/src/gel/core/constants"
	"Gel/src/gel/persistence/repositories"
	"fmt"
	"path/filepath"
)

type IInitService interface {
	Init(path string) (string, error)
}

type InitService struct {
	filesystemRepository repositories.IFilesystemRepository
}

func NewInitService(filesystemRepository repositories.IFilesystemRepository) *InitService {
	return &InitService{
		filesystemRepository,
	}
}

func (initService *InitService) Init(path string) (string, error) {
	base := filepath.Join(path, constants.RepositoryDirName)

	dirs := []string{
		base,
		filepath.Join(base, constants.ObjectsDirName),
		filepath.Join(base, constants.RefsDirName),
	}
	files := []string{
		filepath.Join(base, constants.IndexFileName),
	}

	exists := initService.filesystemRepository.Exists(base)

	for _, dir := range dirs {
		if err := initService.filesystemRepository.MakeDir(dir, constants.DirPermission); err != nil {
			return "", err
		}
	}

	for _, file := range files {
		if err := initService.filesystemRepository.WriteFile(file, []byte{}, false, constants.FilePermission); err != nil {
			return "", err
		}
	}

	if exists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %s", base), nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
