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
	repository repositories.IRepository
}

func NewInitService(repository repositories.IRepository) *InitService {
	return &InitService{
		repository,
	}
}

func (initService *InitService) Init(path string) (string, error) {
	base := filepath.Join(path, constants.RepositoryDirName)

	dirs := []string{
		base,
		filepath.Join(base, constants.ObjectsDirName),
		filepath.Join(base, constants.RefsDirName),
	}
	exists := initService.repository.Exists(base)
	if err := initService.repository.MakeDirRange(dirs, constants.DirPermission); err != nil {
		return err.Error(), err
	}
	if exists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %s", base), nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
