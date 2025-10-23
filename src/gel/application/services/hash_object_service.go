package services

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/persistence/repositories"
	"os"
)

type IHashObjectService interface {
	HashObject(path string, objectType constant.ObjectType, write bool) (string, error)
}

type HashObjectService struct {
	filesystemRepository repositories.IFilesystemRepository
	gelRepository        repositories.IGelRepository
}

func NewHashObjectService(filesystemRepository repositories.IFilesystemRepository,
	gelRepository repositories.IGelRepository) *HashObjectService {
	return &HashObjectService{
		filesystemRepository,
		gelRepository,
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

	writeErr := hashObjectService.filesystemRepository.WriteFile(objectPath, compressedContent, true, constant.FilePermission)

	if writeErr != nil {
		return "", writeErr
	}

	return hash, nil
}
