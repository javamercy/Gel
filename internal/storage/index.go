package storage

import (
	"Gel/internal/workspace"
	"fmt"
	"os"
)

type IndexStorage struct {
	workspaceProvider *workspace.Provider
}

func NewIndexStorage(workspaceProvider *workspace.Provider) *IndexStorage {
	return &IndexStorage{
		workspaceProvider: workspaceProvider,
	}
}

func (i *IndexStorage) Read() ([]byte, error) {
	ws := i.workspaceProvider.GetWorkspace()
	data, err := os.ReadFile(ws.IndexPath)
	if err != nil {
		return nil, fmt.Errorf("error reading index file: %w", err)
	}
	return data, nil
}

func (i *IndexStorage) Write(data []byte) error {
	ws := i.workspaceProvider.GetWorkspace()
	if err := os.WriteFile(ws.IndexPath, data, workspace.FilePermission); err != nil {
		return fmt.Errorf("error writing index file: %w", err)
	}
	return nil
}
