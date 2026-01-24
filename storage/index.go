package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type IndexStorage struct {
	filesystemStorage  *FilesystemStorage
	repositoryProvider *repository.Provider
}

func NewIndexStorage(filesystemStorage *FilesystemStorage, repositoryProvider *repository.Provider) *IndexStorage {
	return &IndexStorage{
		filesystemStorage:  filesystemStorage,
		repositoryProvider: repositoryProvider,
	}
}

func (indexStorage *IndexStorage) Read() ([]byte, error) {
	repo := indexStorage.repositoryProvider.GetRepository()
	data, err := indexStorage.filesystemStorage.ReadFile(repo.IndexPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (indexStorage *IndexStorage) Write(data []byte) error {
	repo := indexStorage.repositoryProvider.GetRepository()
	return indexStorage.filesystemStorage.WriteFile(
		repo.IndexPath,
		data, false,
		constant.GelFilePermission)
}
