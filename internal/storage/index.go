package storage

import (
	workspace2 "Gel/internal/gel/workspace"
	"os"
)

type IndexStorage struct {
	workspaceProvider *workspace2.Provider
}

func NewIndexStorage(workspaceProvider *workspace2.Provider) *IndexStorage {
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
	return os.WriteFile(ws.IndexPath, data, workspace2.FilePermission)
}
