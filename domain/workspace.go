package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotAGelRepository = errors.New("not a Gel repository")
)

type Workspace struct {
	GelDir     string
	ObjectsDir string
	RefsDir    string
	IndexPath  string
	RepoDir    string
	ConfigPath string
}

func NewWorkspace(startPath string) (*Workspace, error) {
	gelDir, err := findGelDir(startPath)
	if err != nil {
		return nil, err
	}
	return &Workspace{
		GelDir:     gelDir,
		ObjectsDir: filepath.Join(gelDir, ObjectsDirName),
		RefsDir:    filepath.Join(gelDir, RefsDirName),
		IndexPath:  filepath.Join(gelDir, IndexFileName),
		RepoDir:    filepath.Dir(gelDir),
		ConfigPath: filepath.Join(gelDir, ConfigFileName),
	}, nil
}

func findGelDir(startPath string) (string, error) {
	currentPath := startPath
	for {
		gelPath := filepath.Join(currentPath, GelDirName)
		info, err := os.Stat(gelPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				parentPath := filepath.Dir(currentPath)
				if parentPath == currentPath {
					break
				}
				currentPath = parentPath
				continue
			}
			return "", fmt.Errorf("failed to stat '%s': %w", gelPath, err)
		}
		if info.IsDir() {
			return gelPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return "", ErrNotAGelRepository
}
