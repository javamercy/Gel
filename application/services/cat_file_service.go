package services

import (
	"Gel/core/helpers"
	"Gel/domain/objects"
	"Gel/persistence/repositories"
	"os"
	"path/filepath"
)

type ICatFileService interface {
	GetObject(hash string) (objects.IObject, error)
	ObjectExists(hash string) bool
}
type CatFileService struct {
	repository repositories.IRepository
}

func NewCatFileService(repository repositories.IRepository) *CatFileService {
	return &CatFileService{
		repository,
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

	data, err := helpers.Decompress(compressedContent)
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
	objectDir := hash[:2]
	objectFile := hash[2:]
	path := filepath.Join(objectDir, objectFile)
	return catFileService.repository.Exists(path)
}
