package storage

import (
	"Gel/internal/workspace"
	"fmt"
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
	objectPath := o.objectPath(hash)
	dir := filepath.Dir(objectPath)
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(objectPath, data, workspace.FilePermission); err != nil {
		return fmt.Errorf("failed to write object '%s': %w", hash, err)
	}
	return nil
}

func (o *ObjectStorage) Read(hash string) ([]byte, error) {
	objectPath := o.objectPath(hash)
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read object '%s': %w", hash, err)
	}
	return data, nil
}

func (o *ObjectStorage) Exists(hash string) (bool, error) {
	objectPath := o.objectPath(hash)
	_, err := os.Stat(objectPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check object '%s' existence: %w", hash, err)
}

func (o *ObjectStorage) objectPath(hash string) string {
	w := o.workspaceProvider.GetWorkspace()
	dir := hash[:2]
	file := hash[2:]
	return filepath.Join(w.ObjectsDir, dir, file)
}
