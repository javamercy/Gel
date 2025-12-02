package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/domain"
	"Gel/src/gel/persistence/repositories"
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
	stringBuilder := strings.Builder{}
	for _, entry := range index.Entries {
		stringBuilder.WriteString(strconv.Itoa(int(entry.Mode)))
		stringBuilder.WriteString(constant.Space)
		stringBuilder.WriteString(entry.Hash)
		stringBuilder.WriteString(constant.Space)
		stringBuilder.WriteString(strconv.Itoa(int(entry.GetStage())))
		stringBuilder.WriteString(constant.Space)
		stringBuilder.WriteString(entry.Path)
		stringBuilder.WriteString(constant.NewLine)

	}
	return stringBuilder.String()
}

func lsFiles(index *domain.Index) string {
	stringBuilder := strings.Builder{}
	for _, entry := range index.Entries {
		stringBuilder.WriteString(entry.Path)
		stringBuilder.WriteString(constant.NewLine)
	}
	return stringBuilder.String()
}
