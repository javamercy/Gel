package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotAGelRepository = fmt.Errorf("not a Gel repository (%s not found)", GelDirName)
)

type Provider struct {
	workspace *Workspace
}

func NewProvider(path string) (*Provider, error) {
	workspace, err := NewWorkspaceFromPath(path)
	if err != nil {
		return nil, err
	}
	return &Provider{
		workspace: workspace,
	}, nil
}

func (provider *Provider) GetWorkspace() *Workspace {
	return provider.workspace
}

type Workspace struct {
	GelDir     string
	ObjectsDir string
	RefsDir    string
	IndexPath  string
	RepoDir    string
	ConfigPath string
}

func NewWorkspaceFromPath(path string) (*Workspace, error) {
	gelDir, err := findGelDir(path)
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
		if err == nil && info.IsDir() {
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
