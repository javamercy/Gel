package vcs

import (
	"Gel/core/constant"
	"Gel/core/util"
	"strings"
)

type AddService struct {
	updateIndexService *UpdateIndexService
	pathResolver       *util.PathResolver
}

func NewAddService(updateIndexService *UpdateIndexService, pathResolver *util.PathResolver) *AddService {
	return &AddService{
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (addService *AddService) Add(pathspecs []string, dryRun bool) (string, error) {
	normalizedPaths, err := addService.pathResolver.Resolve(pathspecs)
	if err != nil {
		return "", err
	}

	// TODO: Git does not print the paths that are already staged
	var result strings.Builder
	if dryRun {
		for _, path := range normalizedPaths {
			result.WriteString(path)
			result.WriteString(constant.NewLineStr)
		}
	}

	err = addService.updateIndexService.UpdateIndex(normalizedPaths, true, false)

	if err != nil {
		return "", err
	}

	return result.String(), nil
}
