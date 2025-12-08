package vcs

import (
	"Gel/core/utilities"
)

type AddService struct {
	updateIndexService *UpdateIndexService
	pathResolver       *utilities.PathResolver
}

func NewAddService(updateIndexService *UpdateIndexService, pathResolver *utilities.PathResolver) *AddService {
	return &AddService{
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (addService *AddService) Add(pathspecs []string, dryRun bool) ([]string, error) {
	normalizedPaths, err := addService.pathResolver.Resolve(pathspecs)
	if err != nil {
		return nil, err
	}

	if dryRun {
		return normalizedPaths, nil
	}

	addPathErr := addService.updateIndexService.UpdateIndex(normalizedPaths, true, false)

	if addPathErr != nil {
		return nil, addPathErr
	}
	return normalizedPaths, nil
}
