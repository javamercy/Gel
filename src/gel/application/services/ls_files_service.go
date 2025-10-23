package services

import (
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/persistence/repositories"
)

type LsFilesOptions struct {
	Stage  bool
	Cached bool
}

type ILsFilesService interface {
	LsFiles() ([]string, error)
}

type LsFilesService struct {
	gelRepository        repositories.IGelRepository
	filesystemRepository repositories.IFilesystemRepository
}

func NewLsFilesService(gelRepository repositories.IGelRepository, filesystemRepository repositories.IFilesystemRepository) *LsFilesService {
	return &LsFilesService{
		gelRepository,
		filesystemRepository,
	}
}

func (lsFilesService *LsFilesService) LsFiles() ([]string, error) {

	indexFilePath, err := lsFilesService.gelRepository.FindIndexFilePath(".")
	if err != nil {
		return nil, err
	}

	indexBytes, err := lsFilesService.filesystemRepository.ReadFile(indexFilePath)
	if err != nil {
		return nil, err
	}

	index, err := serialization.DeserializeIndex(indexBytes)
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(index.Entries))
	for _, entry := range index.Entries {
		files = append(files, entry.Path)
	}
	return files, nil
}
