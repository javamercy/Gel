package storage

import (
	workspace2 "Gel/internal/gel/workspace"
	"os"
	"path/filepath"
)

type ObjectStorage struct {
	workspaceProvider *workspace2.Provider
}

func NewObjectStorage(workspaceProvider *workspace2.Provider) *ObjectStorage {
	return &ObjectStorage{
		workspaceProvider: workspaceProvider,
	}
}

func (o *ObjectStorage) Write(hash string, data []byte) error {
	objectPath := o.objectPath(hash)
	dir := filepath.Dir(objectPath)
	if err := os.MkdirAll(dir, workspace2.DirPermission); err != nil {
		return err
	}
	return os.WriteFile(objectPath, data, workspace2.FilePermission)
}

func (o *ObjectStorage) Read(hash string) ([]byte, error) {
	objectPath := o.objectPath(hash)
	return os.ReadFile(objectPath)
}

func (o *ObjectStorage) Exists(hash string) bool {
	objectPath := o.objectPath(hash)
	_, err := os.Stat(objectPath)
	return err == nil
}

func (o *ObjectStorage) objectPath(hash string) string {
	w := o.workspaceProvider.GetWorkspace()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(w.ObjectsDir, dir, file)
}
