package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/persistence/repositories"
	"fmt"
	"path/filepath"
)

type IInitService interface {
	Init(request *dto.InitRequest) (string, *gelErrors.GelError)
}

type InitService struct {
	filesystemRepository repositories.IFilesystemRepository
}

func NewInitService(filesystemRepository repositories.IFilesystemRepository) *InitService {
	return &InitService{
		filesystemRepository,
	}
}

func (initService *InitService) Init(request *dto.InitRequest) (string, *gelErrors.GelError) {

	validator := validators.NewInitValidator()
	validationResult := validator.Validate(request)
	if !validationResult.IsValid() {
		return "", gelErrors.NewGelError(gelErrors.ExitCodeFatal, validationResult.Error())
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
			return "", gelErrors.NewGelError(gelErrors.ExitCodeFatal, fmt.Sprintf("failed to create directory %s: %v", dir, err))
		}
	}

	if exists {
		return "", gelErrors.NewGelError(gelErrors.ExitCodeWarning, "Reinitialized existing Gel repository")
	}

	return fmt.Sprintf("Initialized empty Gel repository in %s", base), nil
}
