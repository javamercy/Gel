package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
	"Gel/storage"
)

type ObjectService struct {
	objectStorage     storage.IObjectStorage
	filesystemService *FilesystemService
}

func NewObjectService(objectStorage storage.IObjectStorage, filesystemService *FilesystemService) *ObjectService {
	return &ObjectService{
		objectStorage:     objectStorage,
		filesystemService: filesystemService,
	}
}

func (objectService *ObjectService) GetObjectSize(hash string) (uint32, error) {
	compressedContent, err := objectService.objectStorage.Read(hash)
	if err != nil {
		return 0, err
	}

	content, err := encoding.Decompress(compressedContent)
	if err != nil {
		return 0, err
	}

	object, err := domain.DeserializeObject(content)
	if err != nil {
		return 0, err
	}
	return uint32(object.Size()), nil
}

func (objectService *ObjectService) Write(hash string, content []byte) error {
	compressedContent, err := encoding.Compress(content)
	if err != nil {
		return err
	}
	return objectService.objectStorage.Write(hash, compressedContent)
}

func (objectService *ObjectService) Read(hash string) (domain.IObject, error) {
	compressedContent, err := objectService.objectStorage.Read(hash)
	if err != nil {
		return nil, err
	}

	content, err := encoding.Decompress(compressedContent)
	if err != nil {
		return nil, err
	}
	object, err := domain.DeserializeObject(content)
	if err != nil {
		return nil, err
	}
	return object, nil
}
