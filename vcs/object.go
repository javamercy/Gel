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

func (objectService *ObjectService) ReadTreeAndDeserializeEntries(treeHash string) ([]*domain.TreeEntry, error) {
	object, err := objectService.Read(treeHash)
	if err != nil {
		return nil, err
	}

	tree, ok := object.(*domain.Tree)
	if !ok {
		return nil, err
	}

	treeEntries, err := tree.DeserializeTree()
	if err != nil {
		return nil, err
	}

	return treeEntries, nil
}

func (objectService *ObjectService) HashObject(path string) (string, error) {
	data, err := objectService.filesystemService.ReadFile(path)
	if err != nil {
		return "", err
	}
	blob := domain.NewBlob(data)
	content := blob.Serialize()
	return encoding.ComputeHash(content), nil
}
