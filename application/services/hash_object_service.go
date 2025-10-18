package services

import (
	"Gel/core/constants"
	"Gel/core/helpers"
	"Gel/persistence/repositories"
	"os"
)

type IHashObjectService interface {
	HashObject(path string, objectType constants.ObjectType, write bool) (string, error)
}

type HashObjectService struct {
	repository repositories.IRepository
}

func NewHashObjectService(repository repositories.IRepository) *HashObjectService {
	return &HashObjectService{
		repository,
	}
}

func (hashObjectService *HashObjectService) HashObject(path string, objectType constants.ObjectType, write bool) (string, error) {
	fileData, err := hashObjectService.repository.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := helpers.ToObjectContent(objectType, fileData)
	hash := helpers.ComputeHash(content)

	if !write {
		return hash, nil
	}

	compressedContent, err := helpers.Compress(content)
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
