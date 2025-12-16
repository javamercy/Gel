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
	filesystemStorage IFilesystemStorage
}

func NewIndexStorage(filesystemStorage IFilesystemStorage) *IndexStorage {
	return &IndexStorage{
		filesystemStorage: filesystemStorage,
	}
}

func (indexStorage *IndexStorage) Read() (*domain.Index, error) {
	repo := repository.GetRepository()
	data, err := indexStorage.filesystemStorage.ReadFile(repo.IndexPath)
	if err != nil {
		return nil, err
	}

	return domain.DeserializeIndex(data)
}

func (indexStorage *IndexStorage) Write(index *domain.Index) error {
	repo := repository.GetRepository()
	content, err := index.Serialize()
	if err != nil {
		return err
	}
	return indexStorage.filesystemStorage.WriteFile(
		repo.IndexPath,
		content, false,
		constant.GelFilePermission)
}
