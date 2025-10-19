package services

import (
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
	"os"
)

type ICatFileService interface {
	GetObject(hash string) (objects.IObject, error)
	ObjectExists(hash string) bool
}
type CatFileService struct {
	repository        repositories.IRepository
	compressionHelper helpers.ICompressionHelper
}

func NewCatFileService(repository repositories.IRepository, compressionHelper helpers.ICompressionHelper) *CatFileService {
	return &CatFileService{
		repository,
		compressionHelper,
	}
}

func (catFileService *CatFileService) GetObject(hash string) (objects.IObject, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := catFileService.repository.FindObjectPath(hash, cwd)
	compressedContent, err := catFileService.repository.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, err := catFileService.compressionHelper.Decompress(compressedContent)
	if err != nil {
		return nil, err
	}

	object, err := helpers.ToObject(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (catFileService *CatFileService) ObjectExists(hash string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	path, err := catFileService.repository.FindObjectPath(hash, cwd)
	if err != nil {
		return false
	}

	return catFileService.repository.Exists(path)
}
