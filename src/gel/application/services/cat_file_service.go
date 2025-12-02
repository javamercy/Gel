package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
)

type ICatFileService interface {
	GetObject(request *dto.CatFileRequest) (objects.IObject, *gelErrors.GelError)
}

type CatFileService struct {
	filesystemRepository repositories.IFilesystemRepository
	objectRepository     repositories.IObjectRepository
	catFileRules         *rules.CatFileRules
}

func NewCatFileService(
	filesystemRepository repositories.IFilesystemRepository,
	objectRepository repositories.IObjectRepository,
	catFileRules *rules.CatFileRules) *CatFileService {
	return &CatFileService{
		filesystemRepository,
		objectRepository,
		catFileRules,
	}
}

func (catFileService *CatFileService) GetObject(request *dto.CatFileRequest) (objects.IObject, *gelErrors.GelError) {

	validator := validators.NewCatFileValidator()
	validationResult := validator.Validate(request)

	if !validationResult.IsValid() {
		return nil,
			gelErrors.NewGelError(gelErrors.ExitCodeFatal, validationResult.Error())
	}

	err := utilities.RunAll(
		catFileService.catFileRules.ObjectMustExist(request.Hash),
	)

	if err != nil {
		return nil,
			gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	compressedContent, err := catFileService.objectRepository.Read(request.Hash)
	if err != nil {
		return nil,
			gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	content, err := encoding.Decompress(compressedContent)
	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	object, err := objects.Deserialize(content)
	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	return object, nil
}
