package services

import (
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
)

type ICatFileService interface {
	GetObject(hash string) (objects.IObject, error)
}
type CatFileService struct {
	filesystemRepository repositories.IFilesystemRepository
	objectRepository     repositories.IObjectRepository
}

func NewCatFileService(filesystemRepository repositories.IFilesystemRepository, objectRepository repositories.IObjectRepository) *CatFileService {
	return &CatFileService{
		filesystemRepository,
		objectRepository,
	}
}

func (catFileService *CatFileService) GetObject(hash string) (objects.IObject, error) {

	compressedContent, err := catFileService.objectRepository.Read(hash)
	if err != nil {
		return nil, err
	}

	data, err := encoding.Decompress(compressedContent)
	if err != nil {
		return nil, err
	}

	object, err := serialization.DeserializeObject(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}
