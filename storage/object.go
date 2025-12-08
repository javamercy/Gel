package storage

import (
	"Gel/core/constant"
	"Gel/core/context"
	"path/filepath"
)

type IObjectStorage interface {
	Write(hash string, data []byte) error
	Read(hash string) ([]byte, error)
	Exists(hash string) bool
	GetObjectPath(hash string) string
}

type ObjectStorage struct {
	filesystemStorage IFilesystemStorage
}

func NewObjectStorage(filesystemStorage IFilesystemStorage) *ObjectStorage {
	return &ObjectStorage{

		filesystemStorage: filesystemStorage,
	}
}

func (objectStorage *ObjectStorage) Write(hash string, data []byte) error {
	objectPath := objectStorage.GetObjectPath(hash)
	return objectStorage.filesystemStorage.WriteFile(objectPath, data, true, constant.GelFilePermission)
}

func (objectStorage *ObjectStorage) Read(hash string) ([]byte, error) {
	objectPath := objectStorage.GetObjectPath(hash)
	return objectStorage.filesystemStorage.ReadFile(objectPath)
}

func (objectStorage *ObjectStorage) Exists(hash string) bool {
	objectPath := objectStorage.GetObjectPath(hash)
	return objectStorage.filesystemStorage.Exists(objectPath)
}

func (objectStorage *ObjectStorage) GetObjectPath(hash string) string {
	ctx := context.GetContext()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(ctx.ObjectsDir, dir, file)
}
