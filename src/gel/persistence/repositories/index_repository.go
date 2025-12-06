package repositories

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/context"
	"Gel/src/gel/domain"
)

type IIndexRepository interface {
	Read() (*domain.Index, error)
	Write(index *domain.Index) error
	GetEntries() ([]*domain.IndexEntry, error)
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

	return domain.DeserializeIndex(data)
}

func (indexRepository *IndexRepository) Write(index *domain.Index) error {
	ctx := context.GetContext()
	content := index.Serialize()
	return indexRepository.filesystemRepository.WriteFile(
		ctx.IndexPath,
		content, false,
		constant.GelFilePermission)
}

func (indexRepository *IndexRepository) GetEntries() ([]*domain.IndexEntry, error) {
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
