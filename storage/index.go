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
	repository        *repository.Repository
}

func NewIndexStorage(filesystemStorage IFilesystemStorage, repository *repository.Repository) *IndexStorage {
	return &IndexStorage{
		filesystemStorage: filesystemStorage,
		repository:        repository,
	}
}

func (indexStorage *IndexStorage) Read() (*domain.Index, error) {
	data, err := indexStorage.filesystemStorage.ReadFile(indexStorage.repository.IndexPath)
	if err != nil {
		return nil, err
	}

	return domain.DeserializeIndex(data)
}

func (indexStorage *IndexStorage) Write(index *domain.Index) error {
	content, err := index.Serialize()
	if err != nil {
		return err
	}
	return indexStorage.filesystemStorage.WriteFile(
		indexStorage.repository.IndexPath,
		content, false,
		constant.GelFilePermission)
}
