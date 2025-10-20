package services

import (
	"Gel/src/gel/core/constants"
	"Gel/src/gel/core/helpers"
	"Gel/src/gel/persistence/repositories"
	"os"
)

type IHashObjectService interface {
	HashObject(path string, objectType constants.ObjectType, write bool) (string, error)
}

type HashObjectService struct {
	filesystemRepository repositories.IFilesystemRepository
	gelRepository        repositories.IGelRepository
	compressionHelper    helpers.ICompressionHelper
}

func NewHashObjectService(filesystemRepository repositories.IFilesystemRepository,
	gelRepository repositories.IGelRepository, compressionHelper helpers.ICompressionHelper) *HashObjectService {
	return &HashObjectService{
		filesystemRepository,
		gelRepository,
		compressionHelper,
	}
}

func (hashObjectService *HashObjectService) HashObject(path string, objectType constants.ObjectType, write bool) (string, error) {
	fileData, err := hashObjectService.filesystemRepository.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := helpers.SerializeObject(objectType, fileData)
	hash := helpers.ComputeHash(content)

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

	objectPath, err := hashObjectService.gelRepository.FindObjectPath(hash, cwd)
	if err != nil {
		return "", err
	}

	if hashObjectService.filesystemRepository.Exists(objectPath) {
		return hash, nil
	}

	writeErr := hashObjectService.filesystemRepository.WriteFile(objectPath, compressedContent, true, constants.FilePermission)

	if writeErr != nil {
		return "", writeErr
	}

	return hash, nil
}
