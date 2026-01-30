package storage

import (
	"Gel/internal/workspace"
	"path/filepath"
)

type ObjectStorage struct {
	filesystemStorage *FilesystemStorage
	workspaceProvider *workspace.Provider
}

func NewObjectStorage(filesystemStorage *FilesystemStorage, workspaceProvider *workspace.Provider) *ObjectStorage {
	return &ObjectStorage{
		filesystemStorage: filesystemStorage,
		workspaceProvider: workspaceProvider,
	}
}

func (objectStorage *ObjectStorage) Write(hash string, data []byte) error {
	objectPath := objectStorage.GetObjectPath(hash)
	return objectStorage.filesystemStorage.WriteFile(objectPath, data, true, workspace.FilePermission)
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
	ws := objectStorage.workspaceProvider.GetWorkspace()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(ws.ObjectsDir, dir, file)
}
