package repositories

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/context"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/domain"
)

type IIndexRepository interface {
	Read() (*domain.Index, error)
	Write(index *domain.Index) error
	GetAllEntries() ([]*domain.IndexEntry, error)
	AddOrUpdateEntry(entry *domain.IndexEntry) error
	AddOrUpdateEntries(entries []*domain.IndexEntry) error
}

type IndexRepository struct {
	filesystemRepository IFilesystemRepository
}

func NewIndexRepository(filesystemRepository IFilesystemRepository) *IndexRepository {
	return &IndexRepository{
		filesystemRepository: filesystemRepository,
	}
}

func (indexRepository *IndexRepository) Read() (*domain.Index, error) {
	ctx := context.GetContext()
	data, err := indexRepository.filesystemRepository.ReadFile(ctx.IndexPath)
	if err != nil {
		return nil, err
	}

	return serialization.DeserializeIndex(data)
}

func (indexRepository *IndexRepository) Write(index *domain.Index) error {
	ctx := context.GetContext()
	data := serialization.SerializeIndex(index)
	return indexRepository.filesystemRepository.WriteFile(
		ctx.IndexPath,
		data, false,
		constant.GelFilePermission)
}

func (indexRepository *IndexRepository) GetAllEntries() ([]*domain.IndexEntry, error) {
	index, err := indexRepository.Read()
	if err != nil {
		return nil, err
	}
	return index.Entries, nil
}

func (indexRepository *IndexRepository) AddOrUpdateEntry(entry *domain.IndexEntry) error {
	index, err := indexRepository.Read()
	if err != nil {
		return err
	}
	index.AddOrUpdateEntry(entry)
	return indexRepository.Write(index)
}

func (indexRepository *IndexRepository) AddOrUpdateEntries(entries []*domain.IndexEntry) error {
	index, err := indexRepository.Read()
	if err != nil {
		return err
	}
	for _, e := range entries {
		index.AddOrUpdateEntry(e)
	}
	return indexRepository.Write(index)
}
