package repository

import (
	"Gel/core/constant"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrNotAGelRepository = errors.New("not a gel repository")
)

type Repository struct {
	GelDirectory        string
	ObjectsDirectory    string
	RefsDirectory       string
	IndexPath           string
	RepositoryDirectory string
	ConfigPath          string
}

var repository *Repository
var repositoryOnce sync.Once

func GetRepository() *Repository {
	return repository
}

func Initialize() error {
	_, err := initializeRepository()
	if err != nil {
		return ErrNotAGelRepository
	}
	return nil
}

func initializeRepository() (*Repository, error) {
	var err error
	repositoryOnce.Do(func() {
		cwd, e := os.Getwd()
		if e != nil {
			err = e
			return
		}
		gelDirectory, e := findGelDirectory(cwd)
		if e != nil {
			err = e
			return
		}

		repository = &Repository{
			GelDirectory:        gelDirectory,
			ObjectsDirectory:    filepath.Join(gelDirectory, constant.GelObjectsDirectoryName),
			RefsDirectory:       filepath.Join(gelDirectory, constant.GelRefsDirectoryName),
			IndexPath:           filepath.Join(gelDirectory, constant.GelIndexFileName),
			RepositoryDirectory: filepath.Dir(gelDirectory),
			ConfigPath:          filepath.Join(gelDirectory, constant.GelConfigFileName),
		}
	})
	return repository, err
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

	return "", os.ErrNotExist
}
