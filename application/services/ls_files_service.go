package services

import (
	"Gel/application/dto"
	"Gel/core/constant"
	"Gel/core/crossCuttingConcerns/gelErrors"
	"Gel/core/utilities"
	"Gel/domain"
	"Gel/persistence/repositories"
	"strconv"
	"strings"
)

type ILsFilesService interface {
	LsFiles(request *dto.LsFilesRequest) (string, *gelErrors.GelError)
}

type LsFilesService struct {
	indexRepository repositories.IIndexRepository
}

func NewLsFilesService(indexRepository repositories.IIndexRepository) *LsFilesService {
	return &LsFilesService{
		indexRepository,
	}
}

func (lsFilesService *LsFilesService) LsFiles(request *dto.LsFilesRequest) (string, *gelErrors.GelError) {

	index, err := lsFilesService.indexRepository.Read()
	if err != nil {
		return "", gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	if request.Stage {
		return lsFilesWithStage(index), nil
	}

	return lsFiles(index), nil

}

func lsFilesWithStage(index *domain.Index) string {
	var stringBuilder strings.Builder
	for _, entry := range index.Entries {
		stringBuilder.WriteString(utilities.ConvertModeToString(entry.Mode))
		stringBuilder.WriteString(constant.SpaceStr)
		stringBuilder.WriteString(entry.Hash)
		stringBuilder.WriteString(constant.SpaceStr)
		stringBuilder.WriteString(strconv.Itoa(int(entry.GetStage())))
		stringBuilder.WriteString(constant.SpaceStr)
		stringBuilder.WriteString(entry.Path)
		stringBuilder.WriteString(constant.NewLineStr)

	}
	return stringBuilder.String()
}

func lsFiles(index *domain.Index) string {
	var stringBuilder strings.Builder
	for _, entry := range index.Entries {
		stringBuilder.WriteString(entry.Path)
		stringBuilder.WriteString(constant.NewLineStr)
	}
	return stringBuilder.String()
}
