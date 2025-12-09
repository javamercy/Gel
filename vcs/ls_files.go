package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"strconv"
	"strings"
)

type LsFilesService struct {
	indexService *IndexService
}

func NewLsFilesService(indexService *IndexService) *LsFilesService {
	return &LsFilesService{
		indexService: indexService,
	}
}

func (lsFilesService *LsFilesService) LsFiles(stage bool) (string, error) {
	index, err := lsFilesService.indexService.Read()
	if err != nil {
		return "", err
	}

	if stage {
		return lsFilesWithStage(index), nil
	}

	var stringBuilder strings.Builder
	for _, entry := range index.Entries {
		stringBuilder.WriteString(entry.Path)
		stringBuilder.WriteString(constant.NewLineStr)
	}
	return stringBuilder.String(), nil
}

func lsFilesWithStage(index *domain.Index) string {
	var stringBuilder strings.Builder
	for _, entry := range index.Entries {
		stringBuilder.WriteString(domain.ParseFileMode(entry.Mode).String())
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
