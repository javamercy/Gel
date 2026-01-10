package storage

import (
	"Gel/core/constant"
	"Gel/core/repository"
	"Gel/domain"
)

type IIndexStorage interface {
	Read() (*domain.Index, error)
	Write(index *domain.Index) error
}

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

func (indexStorage *IndexStorage) Read() (*domain.Index, error) {
	repo := indexStorage.repositoryProvider.GetRepository()
	data, err := indexStorage.filesystemStorage.ReadFile(repo.IndexPath)
	if err != nil {
		return nil, err
	}

	return domain.DeserializeIndex(data)
}

func (indexStorage *IndexStorage) Write(index *domain.Index) error {
	repo := indexStorage.repositoryProvider.GetRepository()
	content, err := index.Serialize()
	if err != nil {
		return err
	}
	return indexStorage.filesystemStorage.WriteFile(
		repo.IndexPath,
		content, false,
		constant.GelFilePermission)
}
