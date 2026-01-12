package vcs

import (
	"Gel/core/util"
	"fmt"
	"io"
)

type AddService struct {
	updateIndexService *UpdateIndexService
	pathResolver       util.IPathResolver
}

func NewAddService(updateIndexService *UpdateIndexService, pathResolver util.IPathResolver) *AddService {
	return &AddService{
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (addService *AddService) Add(writer io.Writer, pathspecs []string, dryRun bool) error {
	normalizedPaths, err := addService.pathResolver.Resolve(pathspecs)
	if err != nil {
		return err
	}

	// TODO: Git does not print the paths that are already staged
	if dryRun {
		for _, path := range normalizedPaths {
			if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", path)); err != nil {
				return err
			}
		}
		return nil
	}

	return addService.updateIndexService.UpdateIndex(normalizedPaths, true, false)
}
