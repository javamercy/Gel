package internal

import (
	"Gel/internal/workspace"
	"fmt"
	"os"
	"path/filepath"
)

type InitService struct{}

func NewInitService() *InitService {
	return &InitService{}
}

func (i *InitService) Init(path string) (string, error) {
	base := filepath.Join(path, workspace.GelDirName)

	dirs := []string{
		filepath.Join(base, workspace.ObjectsDirName),
		filepath.Join(base, workspace.RefsDirName, workspace.HeadsDirName),
		filepath.Join(base, workspace.RefsDirName, workspace.TagsDirName),
	}

	exists := fileExists(base)

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
			return "", err
		}
	}

	configPath := filepath.Join(base, workspace.ConfigFileName)
	if err := os.WriteFile(configPath, []byte{}, workspace.FilePermission); err != nil {
		return "", err
	}

	headPath := filepath.Join(base, workspace.HeadFileName)
	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(headPath, headContent, workspace.FilePermission); err != nil {
		return "", err
	}

	if exists {
		return "Reinitialized existing Gel repository", nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", base), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
