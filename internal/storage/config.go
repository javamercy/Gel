package storage

import (
	domain2 "Gel/internal/domain"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigStorage struct {
	workspace *domain2.Workspace
}

func NewConfigStorage(workspace *domain2.Workspace) *ConfigStorage {
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
	if err := os.MkdirAll(dir, domain2.DirPermission); err != nil {
		return fmt.Errorf("config: failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(c.workspace.ConfigPath, data, domain2.FilePermission); err != nil {
		return fmt.Errorf("config: error writing config file: %w", err)
	}
	return nil
}
