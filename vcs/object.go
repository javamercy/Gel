package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
	"Gel/storage"
	"fmt"
)

type ObjectService struct {
	objectStorage     *storage.ObjectStorage
	filesystemStorage *storage.FilesystemStorage
}

func NewObjectService(objectStorage *storage.ObjectStorage, filesystemStorage *storage.FilesystemStorage) *ObjectService {
	return &ObjectService{
		objectStorage:     objectStorage,
		filesystemStorage: filesystemStorage,
	}
}

func (o *ObjectService) GetObjectSize(hash string) (uint32, error) {
	compressedData, err := o.objectStorage.Read(hash)
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

func (o *ObjectService) Write(hash string, data []byte) error {
	compressedData, err := encoding.Compress(data)
	if err != nil {
		return err
	}
	return o.objectStorage.Write(hash, compressedData)
}

func (o *ObjectService) Read(hash string) (domain.IObject, error) {
	compressedData, err := o.objectStorage.Read(hash)
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

func (o *ObjectService) ReadTree(hash string) (*domain.Tree, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}

	tree, ok := object.(*domain.Tree)
	if !ok {
		return nil, domain.ErrInvalidObjectType
	}
	return tree, nil
}
func (o *ObjectService) ReadCommit(hash string) (*domain.Commit, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}

	commit, ok := object.(*domain.Commit)
	if !ok {
		return nil, domain.ErrInvalidObjectType
	}
	return commit, nil
}

func (o *ObjectService) ReadTreeAndDeserializeEntries(treeHash string) ([]domain.TreeEntry, error) {
	object, err := o.Read(treeHash)
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

func (o *ObjectService) ComputeHash(path string) (string, error) {
	fileData, err := o.filesystemStorage.ReadFile(path)
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
