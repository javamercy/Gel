package storage

import (
	"Gel/internal/workspace"
	"os"
	"path/filepath"
)

type ObjectStorage struct {
	workspaceProvider *workspace.Provider
}

func NewObjectStorage(workspaceProvider *workspace.Provider) *ObjectStorage {
	return &ObjectStorage{
		workspaceProvider: workspaceProvider,
	}
}

func (o *ObjectStorage) Write(hash string, data []byte) error {
	objectPath := o.GetObjectPath(hash)
	dir := filepath.Dir(objectPath)
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return err
	}
	return os.WriteFile(objectPath, data, workspace.FilePermission)
}

func (o *ObjectStorage) Read(hash string) ([]byte, error) {
	objectPath := o.GetObjectPath(hash)
	return os.ReadFile(objectPath)
}

func (o *ObjectStorage) Exists(hash string) bool {
	objectPath := o.GetObjectPath(hash)
	_, err := os.Stat(objectPath)
	return err == nil
}

func (o *ObjectStorage) GetObjectPath(hash string) string {
	ws := o.workspaceProvider.GetWorkspace()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(ws.ObjectsDir, dir, file)
}
