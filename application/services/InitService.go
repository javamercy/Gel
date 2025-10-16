package services

import (
	"Gel/persistence/repositories"
	"errors"
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

func (initService *InitService) Init(path string) error {
	exists := initService.repository.Exists(path)
	if exists {
		return errors.New("path already exists")
	}

	base := filepath.Join(path, ".gel")
	dirs := []string{
		base,
		filepath.Join(base, "objects"),
		filepath.Join(base, "refs"),
	}

	for _, dir := range dirs {
		if err := initService.repository.MakeDir(dir); err != nil {
			return err
		}
	}
	return nil
}
