package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNotAGelRepository is returned when no .gel directory is found
// in the current directory or any parent.
var ErrNotAGelRepository = errors.New("not a Gel repository")

// Workspace represents a Gel repository's directory structure.
// It contains paths to all critical repository directories and files.
type Workspace struct {
	// GelDir is the path to the .gel directory.
	GelDir string
	// ObjectsDir is the path to the objects directory.
	ObjectsDir string
	// RefsDir is the path to the refs directory.
	RefsDir string
	// IndexPath is the path to the index file.
	IndexPath string
	// RepoDir is the path to the repository root (parent of .gel).
	RepoDir string
	// ConfigPath is the path to the config file.
	ConfigPath string
}

// NewWorkspace finds the .gel directory starting from startPath
// and returns a Workspace with all paths resolved.
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

// findGelDir walks up the directory tree from startPath looking for .gel.
// Returns the path to .gel or ErrNotAGelRepository if not found.
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
