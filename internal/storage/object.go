package storage

import (
	"Gel/domain"
	"Gel/internal/workspace"
	"errors"
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

func (o *ObjectStorage) Write(hash domain.Hash, data []byte) error {
	objectPath, err := o.getObjectPath(hash)
	if err != nil {
		return err
	}
	dir := filepath.Dir(objectPath.String())
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(objectPath.String(), data, workspace.FilePermission); err != nil {
		return fmt.Errorf("failed to write object '%s': %w", hash, err)
	}
	return nil
}

func (o *ObjectStorage) Read(hash domain.Hash) ([]byte, error) {
	objectPath, err := o.getObjectPath(hash)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(objectPath.String())
	if err != nil {
		return nil, fmt.Errorf("failed to read object '%s': %w", hash, err)
	}
	return data, nil
}

func (o *ObjectStorage) Exists(hash domain.Hash) (bool, error) {
	objectPath, err := o.getObjectPath(hash)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(objectPath.String())
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check object '%s' existence: %w", hash, err)
}

func (o *ObjectStorage) getObjectPath(hash domain.Hash) (domain.AbsolutePath, error) {
	ws := o.workspaceProvider.GetWorkspace()
	hexHash := hash.ToHexString()
	dir := hexHash[:2]
	file := hexHash[2:]
	joined := filepath.Join(ws.ObjectsDir, dir, file)

	absPath, err := domain.NewAbsolutePath(joined)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
