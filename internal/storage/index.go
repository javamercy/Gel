package storage

import (
	"Gel/internal/domain"
	"fmt"
	"os"
)

type IndexStorage struct {
	workspace *domain.Workspace
}

func NewIndexStorage(workspace *domain.Workspace) *IndexStorage {
	return &IndexStorage{
		workspace: workspace,
	}
}

func (i *IndexStorage) Read() ([]byte, error) {
	data, err := os.ReadFile(i.workspace.IndexPath)
	if err != nil {
		return nil, fmt.Errorf("error reading index file: %w", err)
	}
	return data, nil
}

func (i *IndexStorage) Write(data []byte) error {
	if err := os.WriteFile(i.workspace.IndexPath, data, domain.FilePermission); err != nil {
		return fmt.Errorf("error writing index file: %w", err)
	}
	return nil
}
