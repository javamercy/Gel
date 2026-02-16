package storage

import (
	"Gel/internal/workspace"
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
	return os.ReadFile(ws.IndexPath)
}

func (i *IndexStorage) Write(data []byte) error {
	ws := i.workspaceProvider.GetWorkspace()
	return os.WriteFile(ws.IndexPath, data, workspace.FilePermission)
}
