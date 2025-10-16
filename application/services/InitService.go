package services

import (
	"Gel/persistence/repositories"
	"fmt"
	"path/filepath"
)

type InitService struct {
	repository repositories.IRepository
}

func NewInitService(repository repositories.IRepository) *InitService {
	return &InitService{
		repository,
	}
}

func (initService *InitService) Init(path string) (string, error) {
	base := filepath.Join(path, ".gel")

	dirs := []string{
		base,
		filepath.Join(base, "objects"),
		filepath.Join(base, "refs"),
	}
	exists := initService.repository.Exists(base)
	if err := initService.repository.MakeDirRange(dirs); err != nil {
		return err.Error(), err
	}
	if exists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %s", base), nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
