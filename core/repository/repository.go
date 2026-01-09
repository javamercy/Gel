package repository

import (
	"Gel/core/constant"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrNotAGelRepository = errors.New(fmt.Sprintf("not a Gel repository (%s not found)", constant.GelRepositoryName))
)

type IRepositoryProvider interface {
	GetRepository() *Repository
	HasRepository() bool
}

type RepositoryProvider struct {
	repository *Repository
}

func NewRepositoryProvider(path string) *RepositoryProvider {
	repository, err := NewRepositoryFromPath(path)
	if err != nil {
		return &RepositoryProvider{
			repository: nil,
		}
	}
	return &RepositoryProvider{
		repository: repository,
	}
}

func (repositoryProvider *RepositoryProvider) GetRepository() *Repository {
	return repositoryProvider.repository
}

func (repositoryProvider *RepositoryProvider) HasRepository() bool {
	return repositoryProvider.repository != nil
}

type Repository struct {
	GelDirectory        string
	ObjectsDirectory    string
	RefsDirectory       string
	IndexPath           string
	RepositoryDirectory string
	ConfigPath          string
}

func NewRepositoryFromPath(path string) (*Repository, error) {
	gelDirectory, err := findGelDirectory(path)
	if err != nil {
		return nil, err
	}
	return &Repository{
		GelDirectory:        gelDirectory,
		ObjectsDirectory:    filepath.Join(gelDirectory, constant.GelObjectsDirectoryName),
		RefsDirectory:       filepath.Join(gelDirectory, constant.GelRefsDirectoryName),
		IndexPath:           filepath.Join(gelDirectory, constant.GelIndexFileName),
		RepositoryDirectory: filepath.Dir(gelDirectory),
		ConfigPath:          filepath.Join(gelDirectory, constant.GelConfigFileName),
	}, nil
}

func findGelDirectory(startPath string) (string, error) {
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
