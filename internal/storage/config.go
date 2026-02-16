package storage

import (
	workspace2 "Gel/internal/gel/workspace"
	"os"
	"path/filepath"
)

type ConfigStorage struct {
	workspaceProvider *workspace2.Provider
}

func NewConfigStorage(workspaceProvider *workspace2.Provider) *ConfigStorage {
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
	if err := os.MkdirAll(dir, workspace2.DirPermission); err != nil {
		return err
	}
	return os.WriteFile(ws.ConfigPath, data, workspace2.FilePermission)
}
