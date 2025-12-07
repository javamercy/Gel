package services

import (
	"Gel/application/dto"
	"Gel/application/rules"
	"Gel/application/validators"
	"Gel/core/crossCuttingConcerns/gelErrors"
	"Gel/core/encoding"
	"Gel/core/utilities"
	"Gel/domain/objects"
	repositories2 "Gel/persistence/repositories"
)

type ICatFileService interface {
	GetObject(request *dto.CatFileRequest) (objects.IObject, *gelErrors.GelError)
}

type CatFileService struct {
	filesystemRepository repositories2.IFilesystemRepository
	objectRepository     repositories2.IObjectRepository
	catFileRules         *rules.CatFileRules
}

func NewCatFileService(
	filesystemRepository repositories2.IFilesystemRepository,
	objectRepository repositories2.IObjectRepository,
	catFileRules *rules.CatFileRules) *CatFileService {
	return &CatFileService{
		filesystemRepository,
		objectRepository,
		catFileRules,
	}
}

func (catFileService *CatFileService) GetObject(request *dto.CatFileRequest) (objects.IObject, *gelErrors.GelError) {

	validator := validators.NewCatFileValidator()
	gelError := validator.Validate(request)
	if gelError != nil {
		return nil, gelError
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

	object, err := objects.DeserializeObject(content)
	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	return object, nil
}
