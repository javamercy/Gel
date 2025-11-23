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
