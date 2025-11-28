package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/persistence/repositories"
	"errors"
	"fmt"
	"path/filepath"
)

type IInitService interface {
	Init(request *dto.InitRequest) (string, error)
}

type InitService struct {
	filesystemRepository repositories.IFilesystemRepository
}

func NewInitService(filesystemRepository repositories.IFilesystemRepository) *InitService {
	return &InitService{
		filesystemRepository,
	}
}

func (initService *InitService) Init(request *dto.InitRequest) (string, error) {

	validator := validators.NewInitValidator()
	validationResult := validator.Validate(request)
	if !validationResult.IsValid() {
		return "", errors.New(validationResult.Error())
	}

	base := filepath.Join(request.Path, constant.GelDirName)

	dirs := []string{
		base,
		filepath.Join(base, constant.GelObjectsDirName),
		filepath.Join(base, constant.GelRefsDirName),
	}

	exists := initService.filesystemRepository.Exists(base)

	for _, dir := range dirs {
		if err := initService.filesystemRepository.MakeDir(dir, constant.GelDirPermission); err != nil {
			return "", err
		}
	}

	if exists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %s", base), nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
