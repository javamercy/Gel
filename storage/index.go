package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
)

type IIndexStorage interface {
	Read() ([]byte, error)
	Write(index []byte) error
}

var _ IIndexStorage = (*IndexStorage)(nil)

type IndexStorage struct {
	filesystemStorage  IFilesystemStorage
	repositoryProvider repository.IRepositoryProvider
}

func NewIndexStorage(filesystemStorage IFilesystemStorage, repositoryProvider repository.IRepositoryProvider) *IndexStorage {
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
