package storage

import (
	"Gel/internal/domain"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigStorage struct {
	workspace *domain.Workspace
}

func NewConfigStorage(workspace *domain.Workspace) *ConfigStorage {
	return &ConfigStorage{
		workspace: workspace,
	}
}

func (c *ConfigStorage) Read() ([]byte, error) {
	data, err := os.ReadFile(c.workspace.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("config: error reading config file: %w", err)
	}
	return data, nil
}

func (c *ConfigStorage) Write(data []byte) error {
	dir := filepath.Dir(c.workspace.ConfigPath)
	if err := os.MkdirAll(dir, domain.DirPermission); err != nil {
		return fmt.Errorf("config: failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(c.workspace.ConfigPath, data, domain.FilePermission); err != nil {
		return fmt.Errorf("config: error writing config file: %w", err)
	}
	return nil
}
