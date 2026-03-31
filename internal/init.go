package internal

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/validate"
	"fmt"
	"os"
	"path/filepath"
)

type InitService struct {
}

func NewInitService() *InitService {
	return &InitService{}
}

func (i *InitService) Init(path string) (string, error) {
	if err := validate.StringMustNotBeEmpty(path); err != nil {
		return "", fmt.Errorf("init: invalid path '%s': %w", path, err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("init: %w", err)
	}

	repoExists, err := core.Exists(absPath)
	if err != nil {
		return "", fmt.Errorf("init: %w", err)
	}

	if repoExists {
		if err := validate.PathMustBeDirectory(absPath); err != nil {
			return "", fmt.Errorf("init: %w", err)
		}
	} else {
		if err := os.MkdirAll(absPath, domain.DirPermission); err != nil {
			return "", fmt.Errorf("init: failed to create directory '%s': %w", absPath, err)
		}
	}

	gelPath := filepath.Join(absPath, domain.GelDirName)
	objectsPath := filepath.Join(gelPath, domain.ObjectsDirName)
	headsPath := filepath.Join(gelPath, domain.RefsDirName, domain.HeadsDirName)
	headPath := filepath.Join(gelPath, domain.HeadFileName)
	configPath := filepath.Join(gelPath, domain.ConfigFileName)

	gelExists, err := core.Exists(gelPath)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(objectsPath, domain.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create objects directory at '%s': %w", objectsPath, err)
	}
	if err := os.MkdirAll(headsPath, domain.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create heads directory at '%s': %w", headsPath, err)
	}

	headExists, err := core.Exists(headPath)
	if err != nil {
		return "", err
	}
	if !headExists {
		headRef := fmt.Sprintf("ref: %s\n", domain.MainRef)
		if err := os.WriteFile(headPath, []byte(headRef), domain.FilePermission); err != nil {
			return "", fmt.Errorf("init: failed to create head file at '%s': %w", headPath, err)
		}
	}

	configExists, err := core.Exists(configPath)
	if err != nil {
		return "", err
	}
	if !configExists {
		if err := os.WriteFile(configPath, nil, domain.FilePermission); err != nil {
			return "", fmt.Errorf("init: failed to create config file at '%s': %w", configPath, err)
		}
	}

	if gelExists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %v", gelPath), nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", gelPath), nil
}
