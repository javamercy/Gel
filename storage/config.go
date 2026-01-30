package storage

import (
	"Gel/internal/workspace"
)

type ConfigStorage struct {
	filesystemStorage *FilesystemStorage
	workspaceProvider *workspace.Provider
}

func NewConfigStorage(filesystemStorage *FilesystemStorage, workspaceProvider *workspace.Provider) *ConfigStorage {
	return &ConfigStorage{
		filesystemStorage: filesystemStorage,
		workspaceProvider: workspaceProvider,
	}
}

func (configStorage *ConfigStorage) Read() ([]byte, error) {
	ws := configStorage.workspaceProvider.GetWorkspace()
	data, err := configStorage.filesystemStorage.ReadFile(ws.ConfigPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (configStorage *ConfigStorage) Write(data []byte) error {
	ws := configStorage.workspaceProvider.GetWorkspace()
	return configStorage.filesystemStorage.WriteFile(
		ws.ConfigPath,
		data,
		true,
		workspace.FilePermission)
}
