package services

import (
	"Gel/core/constants"
	"Gel/core/helpers"
	"Gel/persistence/repositories"
	"os"
	"path/filepath"
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

	content := helpers.PrepareObjectContent(objectType, fileData)
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
	gelDir, err := hashObjectService.repository.FindGelDir(cwd)
	if err != nil {
		return "", err
	}

	objectDir := filepath.Join(gelDir, constants.ObjectsDirName, hash[:2])
	objectPath := filepath.Join(objectDir, hash[2:])

	if err := hashObjectService.repository.MakeDir(objectDir); err != nil {
		return "", err
	}
	if err := hashObjectService.repository.WriteFile(objectPath, compressedContent); err != nil {
		return "", err
	}

	return hash, nil
}
