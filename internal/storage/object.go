package storage

import (
	"Gel/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ObjectStorage struct {
	workspace *domain.Workspace
}

func NewObjectStorage(workspace *domain.Workspace) *ObjectStorage {
	return &ObjectStorage{
		workspace: workspace,
	}
}

func (o *ObjectStorage) Write(hash domain.Hash, data []byte) error {
	objectPath, err := o.getObjectPath(hash)
	if err != nil {
		return err
	}
	dir := filepath.Dir(objectPath.String())
	if err := os.MkdirAll(dir, domain.DirPermission); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(objectPath.String(), data, domain.FilePermission); err != nil {
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
	hexHash := hash.ToHexString()
	dir := hexHash[:2]
	file := hexHash[2:]
	joined := filepath.Join(o.workspace.ObjectsDir, dir, file)
	absPath, err := domain.NewAbsolutePath(joined)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
