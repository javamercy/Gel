package internal

import (
	"Gel/internal/workspace"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type InitService struct {
}

func NewInitService() *InitService {
	return &InitService{}
}

func (i *InitService) InitAndOutput(writer io.Writer, path string) error {
	message, err := i.Init(path)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "%s\n", message)
	return err
}

func (i *InitService) Init(path string) (string, error) {
	repoPath := filepath.Join(path, workspace.GelDirName)

	objectsPath := filepath.Join(repoPath, workspace.ObjectsDirName)
	headsPath := filepath.Join(repoPath, workspace.RefsDirName, workspace.HeadsDirName)
	headPath := filepath.Join(headsPath, workspace.HeadFileName)
	configPath := filepath.Join(repoPath, workspace.ConfigFileName)

	gelExists, err := i.exists(repoPath)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(objectsPath, workspace.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create objects directory at '%s': %w", objectsPath, err)
	}
	if err := os.MkdirAll(headsPath, workspace.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create heads directory at '%s': %w", headsPath, err)
	}

	headExists, err := i.exists(headPath)
	if err != nil {
		return "", err
	}
	if !headExists {
		headRef := fmt.Sprintf("ref: %s\n", workspace.MainRef)
		if err := os.WriteFile(headPath, []byte(headRef), workspace.FilePermission); err != nil {
			return "", fmt.Errorf("init: failed to create head file at '%s': %w", headPath, err)
		}
	}

	configExists, err := i.exists(configPath)
	if err != nil {
		return "", err
	}
	if !configExists {
		if err := os.WriteFile(configPath, nil, workspace.FilePermission); err != nil {
			return "", fmt.Errorf("init: failed to create config file at '%s': %w", configPath, err)
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("init: failed to get absolute path of '%s': %w", path, err)
	}

	if gelExists {
		return fmt.Sprintf("Reinitialized existing Gel repository in %v", absPath), nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", absPath), nil
}
func (i *InitService) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("init: failed to stat '%s': %w", path, err)
	}
	return true, nil
}
