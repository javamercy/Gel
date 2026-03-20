package internal

import (
	"Gel/internal/pathutil"
	"Gel/internal/validate"
	"Gel/internal/workspace"
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
	if err := validate.StringMustNotBeEmpty(path); err != nil {
		return "", fmt.Errorf("init: invalid path '%s': %w", path, err)
	}
	if err := validate.PathMustExist(path); err != nil {
		return "", fmt.Errorf("init: invalid path '%s': %w", path, err)
	}

	repoPath := filepath.Join(path, workspace.GelDirName)
	objectsPath := filepath.Join(repoPath, workspace.ObjectsDirName)
	headsPath := filepath.Join(repoPath, workspace.RefsDirName, workspace.HeadsDirName)
	headPath := filepath.Join(repoPath, workspace.HeadFileName)
	configPath := filepath.Join(repoPath, workspace.ConfigFileName)

	gelExists, err := pathutil.Exists(repoPath)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(objectsPath, workspace.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create objects directory at '%s': %w", objectsPath, err)
	}
	if err := os.MkdirAll(headsPath, workspace.DirPermission); err != nil {
		return "", fmt.Errorf("init: failed to create heads directory at '%s': %w", headsPath, err)
	}

	headExists, err := pathutil.Exists(headPath)
	if err != nil {
		return "", err
	}
	if !headExists {
		headRef := fmt.Sprintf("ref: %s\n", workspace.MainRef)
		if err := os.WriteFile(headPath, []byte(headRef), workspace.FilePermission); err != nil {
			return "", fmt.Errorf("init: failed to create head file at '%s': %w", headPath, err)
		}
	}

	configExists, err := pathutil.Exists(configPath)
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
