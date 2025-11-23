package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
	"fmt"
)

type ICatFileService interface {
	GetObject(request *dto.CatFileRequest) (objects.IObject, error)
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

func (catFileService *CatFileService) GetObject(request *dto.CatFileRequest) (objects.IObject, error) {

	validator := validators.NewCatFileValidator()
	if err := validator.Validate(request); err != nil {
		return nil, err
	}

	err := utilities.RunAll(
		catFileService.catFileRules.ObjectMustExist(request.Hash),
	)

	if err != nil {
		return nil, err
	}

	compressedContent, err := catFileService.objectRepository.Read(request.Hash)
	if err != nil {
		return nil, err
	}

	data, err := encoding.Decompress(compressedContent)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress object: %v", err)
	}

	object, err := serialization.DeserializeObject(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize object: %v", err)
	}

	return object, nil
}
