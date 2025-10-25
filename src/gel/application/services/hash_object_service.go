package services

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/persistence/repositories"
)

type HashObjectRequest struct {
	Paths      []string
	ObjectType constant.ObjectType
	Write      bool
}

type IHashObjectService interface {
	HashObject(request HashObjectRequest) (map[string]string, error)
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

func (hashObjectService *HashObjectService) HashObject(request HashObjectRequest) (map[string]string, error) {

	hashMap, contentMap, err := hashObjectService.hashObjects(request.Paths, request.ObjectType)

	if err != nil {
		return nil, err
	}
	if !request.Write {
		return hashMap, nil
	}

	err = hashObjectService.write(contentMap)
	if err != nil {
		return nil, err
	}

	return hashMap, nil
}

func (hashObjectService *HashObjectService) hashObjects(paths []string, objectType constant.ObjectType) (map[string]string, map[string][]byte, error) {
	hashMap := make(map[string]string)
	contentMap := make(map[string][]byte)
	for _, path := range paths {
		fileData, err := hashObjectService.filesystemRepository.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		content := serialization.SerializeObject(objectType, fileData)
		hash := encoding.ComputeHash(content)
		hashMap[path] = hash
		contentMap[hash] = content
	}
	return hashMap, contentMap, nil
}

func (hashObjectService *HashObjectService) write(contentMap map[string][]byte) error {
	for hash, content := range contentMap {
		compressedContent, err := encoding.Compress(content)
		if err != nil {
			return err
		}

		writeErr := hashObjectService.objectRepository.Write(hash, compressedContent)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
