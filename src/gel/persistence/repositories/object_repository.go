package repositories

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/context"
	"path/filepath"
)

type IObjectRepository interface {
	Write(hash string, data []byte) error
	Read(hash string) ([]byte, error)
	Exists(hash string) bool
	GetObjectPath(hash string) string
}

type ObjectRepository struct {
	filesystemRepository IFilesystemRepository
}

func NewObjectRepository(filesystemRepository IFilesystemRepository) *ObjectRepository {
	return &ObjectRepository{

		filesystemRepository: filesystemRepository,
	}
}

func (objectRepository *ObjectRepository) Write(hash string, data []byte) error {
	objectPath := objectRepository.GetObjectPath(hash)
	return objectRepository.filesystemRepository.WriteFile(objectPath, data, true, constant.FilePermission)
}

func (objectRepository *ObjectRepository) Read(hash string) ([]byte, error) {
	objectPath := objectRepository.GetObjectPath(hash)
	return objectRepository.filesystemRepository.ReadFile(objectPath)
}

func (objectRepository *ObjectRepository) Exists(hash string) bool {
	objectPath := objectRepository.GetObjectPath(hash)
	return objectRepository.filesystemRepository.Exists(objectPath)
}

func (objectRepository *ObjectRepository) GetObjectPath(hash string) string {
	ctx := context.GetContext()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(ctx.ObjectsDir, dir, file)
}
