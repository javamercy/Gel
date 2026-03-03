package storage

import (
	"Gel/internal/workspace"
	"fmt"
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
	data, err := os.ReadFile(ws.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return data, nil
}

func (c *ConfigStorage) Write(data []byte) error {
	ws := c.workspaceProvider.GetWorkspace()
	dir := filepath.Dir(ws.ConfigPath)
	if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(ws.ConfigPath, data, workspace.FilePermission); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}
