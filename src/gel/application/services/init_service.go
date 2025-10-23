package services

import (
	"Gel/src/gel/core/constant"
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
	base := filepath.Join(path, constant.GelDirName)

	dirs := []string{
		base,
		filepath.Join(base, constant.ObjectsDirName),
		filepath.Join(base, constant.RefsDirName),
	}

	exists := initService.filesystemRepository.Exists(base)

	for _, dir := range dirs {
		if err := initService.filesystemRepository.MakeDir(dir, constant.DirPermission); err != nil {
			return "", err
		}
	}

	if exists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %s", base), nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
