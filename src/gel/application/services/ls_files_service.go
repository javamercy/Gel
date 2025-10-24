package services

import (
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
	indexRepository repositories.IIndexRepository
}

func NewLsFilesService(indexRepository repositories.IIndexRepository) *LsFilesService {
	return &LsFilesService{
		indexRepository,
	}
}

func (lsFilesService *LsFilesService) LsFiles() ([]string, error) {

	index, err := lsFilesService.indexRepository.Read()
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(index.Entries))
	for _, entry := range index.Entries {
		files = append(files, entry.Path)
	}

	return files, nil
}
