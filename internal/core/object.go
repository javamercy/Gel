package core

import (
	domain2 "Gel/internal/domain"
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

func (o *ObjectService) GetObjectSize(hash domain2.Hash) (uint32, error) {
	compressedData, err := o.objectStorage.Read(hash)
	if err != nil {
		return 0, err
	}

	data, err := Decompress(compressedData)
	if err != nil {
		return 0, err
	}

	object, err := domain2.DeserializeObject(data)
	if err != nil {
		return 0, err
	}
	return uint32(object.Size()), nil
}

func (o *ObjectService) Write(hash domain2.Hash, data []byte) error {
	compressedData, err := Compress(data)
	if err != nil {
		return err
	}
	return o.objectStorage.Write(hash, compressedData)
}

func (o *ObjectService) Read(hash domain2.Hash) (domain2.Object, error) {
	compressedData, err := o.objectStorage.Read(hash)
	if err != nil {
		return nil, err
	}

	data, err := Decompress(compressedData)
	if err != nil {
		return nil, err
	}

	object, err := domain2.DeserializeObject(data)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (o *ObjectService) ReadBlob(hash domain2.Hash) (*domain2.Blob, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}
	blob, ok := object.(*domain2.Blob)
	if !ok {
		return nil, domain2.ErrInvalidObjectType
	}
	return blob, nil
}

func (o *ObjectService) ReadTree(hash domain2.Hash) (*domain2.Tree, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}

	tree, ok := object.(*domain2.Tree)
	if !ok {
		return nil, domain2.ErrInvalidObjectType
	}
	return tree, nil
}

func (o *ObjectService) ReadCommit(hash domain2.Hash) (*domain2.Commit, error) {
	object, err := o.Read(hash)
	if err != nil {
		return nil, err
	}

	commit, ok := object.(*domain2.Commit)
	if !ok {
		return nil, domain2.ErrInvalidObjectType
	}
	return commit, nil
}

func (o *ObjectService) ReadTreeAndDeserializeEntries(treeHash domain2.Hash) ([]domain2.TreeEntry, error) {
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

func (o *ObjectService) Exists(hash domain2.Hash) (bool, error) {
	return o.objectStorage.Exists(hash)
}

func (o *ObjectService) ComputeObjectHash(path domain2.AbsolutePath) (domain2.Hash, []byte, error) {
	data, err := os.ReadFile(path.String())
	if err != nil {
		return domain2.Hash{}, nil, fmt.Errorf("failed to read file at '%s': %w", path, err)
	}

	blob := domain2.NewBlob(data)
	serializedData := blob.Serialize()
	hexHash := ComputeSHA256(serializedData)
	hash, err := domain2.NewHash(hexHash)
	if err != nil {
		return domain2.Hash{}, nil, err
	}
	return hash, serializedData, nil
}
