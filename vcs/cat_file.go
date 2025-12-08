package vcs

import (
	"Gel/domain"
)

type CatFileService struct {
	objectService *ObjectService
}

func NewCatFileService(objectService *ObjectService) *CatFileService {
	return &CatFileService{
		objectService: objectService,
	}
}

func (catFileService *CatFileService) CatFile(hash string) (domain.IObject, error) {
	return catFileService.objectService.Read(hash)
}
