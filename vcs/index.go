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
	indexStorage storage.IIndexStorage
}

func NewIndexService(indexStorage storage.IIndexStorage) *IndexService {
	return &IndexService{
		indexStorage: indexStorage,
	}
}

func (indexService *IndexService) Read() (*domain.Index, error) {
	index, err := indexService.indexStorage.Read()
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrIndexNotFound
	}
	return index, err
}

func (indexService *IndexService) Write(index *domain.Index) error {
	return indexService.indexStorage.Write(index)
}

func (indexService *IndexService) GetEntries() ([]*domain.IndexEntry, error) {
	index, err := indexService.Read()
	if err != nil {
		return nil, err
	}
	return index.Entries, nil
}

func (indexService *IndexService) AddOrUpdateEntry(entry *domain.IndexEntry) error {
	index, err := indexService.Read()
	if err != nil {
		return err
	}
	index.AddOrUpdateEntry(entry)
	return indexService.indexStorage.Write(index)
}

func (indexService *IndexService) AddOrUpdateEntries(entries []*domain.IndexEntry) error {
	index, err := indexService.Read()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		index.AddOrUpdateEntry(entry)
	}
	return indexService.indexStorage.Write(index)
}
