package repository

import (
	"Gel/core/constant"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotAGelRepository = fmt.Errorf("not a Gel repository (%s not found)", constant.GelRepositoryName)
)

type IRepositoryProvider interface {
	GetRepository() *Repository
}

// Compile-time interface check
var _ IRepositoryProvider = (*RepositoryProvider)(nil)

type RepositoryProvider struct {
	repository *Repository
}

func NewRepositoryProvider(path string) (*RepositoryProvider, error) {
	repository, err := NewRepositoryFromPath(path)
	if err != nil {
		return nil, err
	}
	return &RepositoryProvider{
		repository: repository,
	}, nil
}

func (repositoryProvider *RepositoryProvider) GetRepository() *Repository {
	return repositoryProvider.repository
}

type Repository struct {
	GelDir        string
	ObjectsDir    string
	RefsDir       string
	IndexPath     string
	RepositoryDir string
	ConfigPath    string
}

func NewRepositoryFromPath(path string) (*Repository, error) {
	gelDir, err := findGelDir(path)
	if err != nil {
		return nil, err
	}
	return &Repository{
		GelDir:        gelDir,
		ObjectsDir:    filepath.Join(gelDir, constant.GelObjectsDirName),
		RefsDir:       filepath.Join(gelDir, constant.GelRefsDirName),
		IndexPath:     filepath.Join(gelDir, constant.GelIndexFileName),
		RepositoryDir: filepath.Dir(gelDir),
		ConfigPath:    filepath.Join(gelDir, constant.GelConfigFileName),
	}, nil
}

func findGelDir(startPath string) (string, error) {
	currentPath := startPath
	for {
		gelPath := filepath.Join(currentPath, constant.GelRepositoryName)
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
