package services

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/persistence/repositories"
)

type IHashObjectService interface {
	HashObject(path string, objectType constant.ObjectType, write bool) (string, error)
}

type HashObjectService struct {
	filesystemRepository repositories.IFilesystemRepository
	objectRepository     repositories.IObjectRepository
}

func NewHashObjectService(filesystemRepository repositories.IFilesystemRepository,
	objectRepository repositories.IObjectRepository) *HashObjectService {
	return &HashObjectService{
		filesystemRepository,
		objectRepository,
	}
}

func (hashObjectService *HashObjectService) HashObject(path string, objectType constant.ObjectType, write bool) (string, error) {

	fileData, err := hashObjectService.filesystemRepository.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := serialization.SerializeObject(objectType, fileData)
	hash := encoding.ComputeHash(content)

	if !write {
		return hash, nil
	}

	compressedContent, err := encoding.Compress(content)
	if err != nil {
		return "", err
	}

	writeErr := hashObjectService.objectRepository.Write(hash, compressedContent)
	if writeErr != nil {
		return "", writeErr
	}

	return hash, nil
}
