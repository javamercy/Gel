package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
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
	repository        *repository.Repository
}

func NewObjectStorage(filesystemStorage IFilesystemStorage, repository *repository.Repository) *ObjectStorage {
	return &ObjectStorage{
		filesystemStorage: filesystemStorage,
		repository:        repository,
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
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(objectStorage.repository.ObjectsDirectory, dir, file)
}
