package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/persistence/repositories"
)

type IAddService interface {
	Add(request *dto.AddRequest) *dto.AddResponse
}

type AddService struct {
	HashObjectService  IHashObjectService
	updateIndexService IUpdateIndexService
	indexRepository    repositories.IIndexRepository
}

func NewAddService(hashObjectService IHashObjectService, updateIndexService IUpdateIndexService, indexRepository repositories.IIndexRepository) *AddService {
	return &AddService{
		HashObjectService:  hashObjectService,
		updateIndexService: updateIndexService,
		indexRepository:    indexRepository,
	}
}

func (s *AddService) Add(request *dto.AddRequest) *dto.AddResponse {
	// 1. Validate request
	// 2. Run rules
	// 3. Resolve paths (patterns, directories, files)
	// 4. Apply filters
	// 5. For each file:
	//    - Hash and store object
	//    - Update index entry
	// 6. Return AddResponse
	return dto.NewAddResponse(nil, nil)
}
