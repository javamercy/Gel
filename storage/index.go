package storage

import (
	"Gel/core/constant"
	"Gel/core/context"
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
	ctx := context.GetContext()
	data, err := indexStorage.filesystemStorage.ReadFile(ctx.IndexPath)
	if err != nil {
		return nil, err
	}

	return domain.DeserializeIndex(data)
}

func (indexStorage *IndexStorage) Write(index *domain.Index) error {
	ctx := context.GetContext()
	content := index.Serialize()
	return indexStorage.filesystemStorage.WriteFile(
		ctx.IndexPath,
		content, false,
		constant.GelFilePermission)
}
