package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ErrNotAGelRepository is returned when repository discovery reaches the
	// filesystem root without finding a .gel directory.
	ErrNotAGelRepository = errors.New("not a Gel repository")

	// ErrInvalidWorkspacePath is returned when workspace discovery receives an
	// invalid starting path.
	ErrInvalidWorkspacePath = errors.New("invalid workspace path")

	// ErrInvalidGelRepository is returned when .gel exists but cannot be used as
	// a repository metadata directory.
	ErrInvalidGelRepository = errors.New("invalid Gel repository")
)

// Workspace represents the absolute, OS-native paths for a Gel repository.
//
// It is a layout descriptor only: constructing a Workspace discovers the
// repository root and derives standard metadata paths, but it does not verify
// that every metadata file exists or is readable.
type Workspace struct {
	// RepoDir is the repository root directory, which is the parent of GelDir.
	RepoDir AbsolutePath

	// GelDir is the .gel repository metadata directory.
	GelDir AbsolutePath

	// ObjectsDir is the .gel/objects object storage directory.
	ObjectsDir AbsolutePath

	// RefsDir is the .gel/refs directory.
	RefsDir AbsolutePath

	// HeadsDir is the .gel/refs/heads directory.
	HeadsDir AbsolutePath

	// HeadPath is the .gel/HEAD symbolic reference file path.
	HeadPath AbsolutePath

	// IndexPath is the .gel/index staging-area file path.
	IndexPath AbsolutePath

	// ConfigPath is the .gel/config.toml repository config file path.
	ConfigPath AbsolutePath
}

// NewWorkspace searches upward from startPath for .gel and returns a Workspace
// with absolute repository paths.
//
// startPath may be relative or absolute. The returned Workspace contains paths
// derived from the first .gel directory found at startPath or one of its
// parents. If a .gel path exists but is not a directory, discovery fails with
// ErrInvalidGelRepository instead of continuing to a parent repository.
func NewWorkspace(startPath string) (*Workspace, error) {
	if startPath == "" {
		return nil, fmt.Errorf("%w: empty path", ErrInvalidWorkspacePath)
	}

	absStartPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf("workspace: resolve start path %q: %w", startPath, err)
	}

	gelDir, err := findGelDir(absStartPath)
	if err != nil {
		return nil, err
	}
	return newWorkspaceFromGelDir(gelDir)
}

// newWorkspaceFromGelDir derives all standard repository paths from gelDir.
// gelDir must be an absolute path to an existing .gel directory.
func newWorkspaceFromGelDir(gelDir string) (*Workspace, error) {
	repoDir, err := newWorkspaceAbsolutePath(filepath.Dir(gelDir))
	if err != nil {
		return nil, err
	}

	gelPath, err := newWorkspaceAbsolutePath(gelDir)
	if err != nil {
		return nil, err
	}

	objectsDir, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, ObjectsDirName))
	if err != nil {
		return nil, err
	}

	refsDir, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, RefsDirName))
	if err != nil {
		return nil, err
	}

	headsDir, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, RefsDirName, HeadsDirName))
	if err != nil {
		return nil, err
	}

	headPath, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, HeadFileName))
	if err != nil {
		return nil, err
	}

	indexPath, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, IndexFileName))
	if err != nil {
		return nil, err
	}

	configPath, err := newWorkspaceAbsolutePath(filepath.Join(gelDir, ConfigFileName))
	if err != nil {
		return nil, err
	}
	return &Workspace{
		RepoDir:    repoDir,
		GelDir:     gelPath,
		ObjectsDir: objectsDir,
		RefsDir:    refsDir,
		HeadsDir:   headsDir,
		HeadPath:   headPath,
		IndexPath:  indexPath,
		ConfigPath: configPath,
	}, nil
}

func newWorkspaceAbsolutePath(path string) (AbsolutePath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return AbsolutePath{}, fmt.Errorf("workspace: derive path %q: %w", path, err)
	}
	return absPath, nil
}

// findGelDir walks upward from startPath and returns the first .gel directory.
// startPath must be absolute. If no .gel directory exists before the filesystem
// root, it returns ErrNotAGelRepository.
func findGelDir(startPath string) (string, error) {
	currentPath := startPath
	for {
		gelPath := filepath.Join(currentPath, GelDirName)
		info, err := os.Stat(gelPath)
		if err == nil {
			if !info.IsDir() {
				return "", fmt.Errorf("%w: %q is not a directory", ErrInvalidGelRepository, gelPath)
			}
			return gelPath, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("workspace: stat %q: %w", gelPath, err)
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", ErrNotAGelRepository
		}
		currentPath = parentPath
	}
}
