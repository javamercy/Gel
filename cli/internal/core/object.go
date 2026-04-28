package core

import (
	"Gel/internal/domain"
	"Gel/internal/storage"
	"fmt"
	"os"
)

type ObjectService struct {
	objectStorage *storage.ObjectStorage
}

func NewObjectService(objectStorage *storage.ObjectStorage) *ObjectService {
	return &ObjectService{
		objectStorage: objectStorage,
	}
}

func (o *ObjectService) GetObjectSize(hash domain.Hash) (uint32, error) {
	compressedData, err := o.objectStorage.Read(hash)
	if err != nil {
		return 0, err
	}

	data, err := Decompress(compressedData)
	if err != nil {
		return 0, err
	}

	object, err := domain.DeserializeObject(data)
	if err != nil {
		return 0, err
	}
	return uint32(object.Size()), nil
}

func (o *ObjectService) Write(hash domain.Hash, data []byte) error {
	compressedData, err := Compress(data)
	if err != nil {
		return err
	}
	return o.objectStorage.Write(hash, compressedData)
}

func (o *ObjectService) Read(hash domain.Hash) (domain.Object, error) {
	compressedData, err := o.objectStorage.Read(hash)
	if err != nil {
		return nil, err
	}

	data, err := Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	object, err := domain.DeserializeObject(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (o *ObjectService) ReadBlob(hash domain.Hash) (*domain.Blob, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}
	blob, ok := object.(*domain.Blob)
	if !ok {
		return nil, domain.ErrInvalidObjectType
	}
	return blob, nil
}

func (o *ObjectService) ReadTree(hash domain.Hash) (*domain.Tree, error) {
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

func (o *ObjectService) ReadCommit(hash domain.Hash) (*domain.Commit, error) {
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

func (o *ObjectService) ReadTreeAndDeserializeEntries(treeHash domain.Hash) ([]domain.TreeEntry, error) {
	tree, err := o.ReadTree(treeHash)
	if err != nil {
		return nil, err
	}

	treeEntries, err := tree.Deserialize()
	if err != nil {
		return nil, err
	}
	return treeEntries, nil
}

func (o *ObjectService) Exists(hash domain.Hash) (bool, error) {
	return o.objectStorage.Exists(hash)
}

func (o *ObjectService) ComputeObjectHash(path domain.AbsolutePath) (domain.Hash, []byte, error) {
	data, err := os.ReadFile(path.String())
	if err != nil {
		return domain.Hash{}, nil, fmt.Errorf("failed to read file at '%s': %w", path, err)
	}

	blob := domain.NewBlob(data)
	serializedData := blob.Serialize()
	hexHash := ComputeSHA256(serializedData)
	hash, err := domain.NewHash(hexHash)
	if err != nil {
		return domain.Hash{}, nil, err
	}
	return hash, serializedData, nil
}
