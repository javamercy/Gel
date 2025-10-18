package services

import (
	"Gel/src/gel/core/constants"
	helpers2 "Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"
	"os"
)

type IHashObjectService interface {
	HashObject(path string, objectType constants.ObjectType, write bool) (string, error)
}

type HashObjectService struct {
	repository        repositories.IRepository
	compressionHelper helpers2.ICompressionHelper
}

func NewHashObjectService(repository repositories.IRepository, compressionHelper helpers2.ICompressionHelper) *HashObjectService {
	return &HashObjectService{
		repository,
		compressionHelper,
	}
}

func (hashObjectService *HashObjectService) HashObject(path string, objectType constants.ObjectType, write bool) (string, error) {
	fileData, err := hashObjectService.repository.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := helpers2.ToObjectContent(objectType, fileData)
	hash := helpers2.ComputeHash(content)

	if !write {
		return hash, nil
	}

	compressedContent, err := hashObjectService.compressionHelper.Compress(content)
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	objectPath, err := hashObjectService.repository.FindObjectPath(hash, cwd)
	if err != nil {
		return "", err
	}

	if hashObjectService.repository.Exists(objectPath) {
		return hash, nil
	}

	writeErr := hashObjectService.repository.WriteFile(objectPath, compressedContent, true, constants.FilePermission)

	if writeErr != nil {
		return "", writeErr
	}

	return hash, nil
}
