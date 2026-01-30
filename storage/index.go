package storage

import (
	"Gel/internal/workspace"
)

type IndexStorage struct {
	filesystemStorage *FilesystemStorage
	workspaceProvider *workspace.Provider
}

func NewIndexStorage(filesystemStorage *FilesystemStorage, workspaceProvider *workspace.Provider) *IndexStorage {
	return &IndexStorage{
		filesystemStorage: filesystemStorage,
		workspaceProvider: workspaceProvider,
	}
}

func (indexStorage *IndexStorage) Read() ([]byte, error) {
	ws := indexStorage.workspaceProvider.GetWorkspace()
	data, err := indexStorage.filesystemStorage.ReadFile(ws.IndexPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (indexStorage *IndexStorage) Write(data []byte) error {
	ws := indexStorage.workspaceProvider.GetWorkspace()
	return indexStorage.filesystemStorage.WriteFile(
		ws.IndexPath,
		data, false,
		workspace.FilePermission)
}
