package rules

import (
	"Gel/src/gel/persistence/repositories"
	"errors"
	"fmt"
)

type UpdateIndexRules struct {
	fileSystemRepository repositories.IFilesystemRepository
}

func NewUpdateIndexRules(fileSystemRepository repositories.IFilesystemRepository) *UpdateIndexRules {
	return &UpdateIndexRules{
		fileSystemRepository,
	}
}

func (updateIndexRules *UpdateIndexRules) AllPathsMustExist(paths []string) error {
	for _, path := range paths {
		exists := updateIndexRules.fileSystemRepository.Exists(path)
		if !exists {
			return errors.New(fmt.Sprintf("Path does not exist: %s", path))
		}
	}
	return nil
}

func (updateIndexRules *UpdateIndexRules) NoDuplicatePaths(paths []string) error {
	pathSet := make(map[string]bool)
	for _, path := range paths {
		if _, exists := pathSet[path]; exists {
			return errors.New(fmt.Sprintf("Duplicate path found: %s", path))
		}
		pathSet[path] = true
	}
	return nil
}

func (updateIndexRules *UpdateIndexRules) PathsMustBeFiles(paths []string) error {
	for _, path := range paths {
		fileInfo, err := updateIndexRules.fileSystemRepository.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("Path is not a file: %s", path))
		}
	}
	return nil
}
