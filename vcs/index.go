package vcs

import (
	"Gel/domain"
	"Gel/storage"
	"errors"
	"io/fs"
)

var (
	ErrIndexNotFound = errors.New("index not found")
)

type IndexService struct {
	indexStorage *storage.IndexStorage
}

func NewIndexService(indexStorage *storage.IndexStorage) *IndexService {
	return &IndexService{
		indexStorage: indexStorage,
	}
}

func (indexService *IndexService) Read() (*domain.Index, error) {
	data, err := indexService.indexStorage.Read()
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrIndexNotFound
	}
	return domain.DeserializeIndex(data)
}

func (indexService *IndexService) Write(index *domain.Index) error {
	serializedData, err := index.Serialize()
	if err != nil {
		return err
	}
	return indexService.indexStorage.Write(serializedData)
}

func (indexService *IndexService) GetEntries() ([]*domain.IndexEntry, error) {
	index, err := indexService.Read()
	if err != nil {
		return nil, err
	}
	return index.Entries, nil
}
