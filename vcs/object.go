package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
	"Gel/storage"
	"fmt"
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
	compressedData, err := objectService.objectStorage.Read(hash)
	if err != nil {
		return 0, err
	}

	data, err := encoding.Decompress(compressedData)
	if err != nil {
		return 0, err
	}

	object, err := domain.DeserializeObject(data)
	if err != nil {
		return 0, err
	}
	return uint32(object.Size()), nil
}

func (objectService *ObjectService) Write(hash string, data []byte) error {
	compressedData, err := encoding.Compress(data)
	if err != nil {
		return err
	}
	return objectService.objectStorage.Write(hash, compressedData)
}

func (objectService *ObjectService) Read(hash string) (domain.IObject, error) {
	compressedData, err := objectService.objectStorage.Read(hash)
	if err != nil {
		return nil, err
	}

	data, err := encoding.Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	object, err := domain.DeserializeObject(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (objectService *ObjectService) ReadTreeAndDeserializeEntries(treeHash string) ([]domain.TreeEntry, error) {
	object, err := objectService.Read(treeHash)
	if err != nil {
		return nil, err
	}

	tree, ok := object.(*domain.Tree)
	if !ok {
		return nil, fmt.Errorf("expected tree object, got %s", object.Type())
	}

	treeEntries, err := tree.Deserialize()
	if err != nil {
		return nil, err
	}
	return treeEntries, nil
}

func (objectService *ObjectService) ComputeHash(path string) (string, error) {
	fileData, err := objectService.filesystemService.ReadFile(path)
	if err != nil {
		return "", err
	}

	blob, err := domain.NewBlob(fileData)
	if err != nil {
		return "", err
	}

	serializedData := blob.Serialize()
	return encoding.ComputeSha256(serializedData), nil
}
