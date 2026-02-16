package storage

import (
	"Gel/internal/workspace"
	"os"
	"path/filepath"
)

type ConfigStorage struct {
	workspaceProvider *workspace.Provider
}

func NewConfigStorage(workspaceProvider *workspace.Provider) *ConfigStorage {
	return &ConfigStorage{
		workspaceProvider: workspaceProvider,
	}
}

func (c *ConfigStorage) Read() ([]byte, error) {
	ws := c.workspaceProvider.GetWorkspace()
	return os.ReadFile(ws.ConfigPath)
}

func (c *ConfigStorage) Write(data []byte) error {
	ws := c.workspaceProvider.GetWorkspace()
	dir := filepath.Dir(ws.ConfigPath)
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return err
	}
	return os.WriteFile(ws.ConfigPath, data, workspace.FilePermission)
}
