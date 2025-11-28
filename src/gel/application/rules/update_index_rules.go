package rules

import (
	"Gel/src/gel/persistence/repositories"
	"errors"
	"fmt"
)

type UpdateIndexRules struct {
	indexRepository repositories.IIndexRepository
}

func NewUpdateIndexRules(indexRepository repositories.IIndexRepository) *UpdateIndexRules {
	return &UpdateIndexRules{
		indexRepository,
	}
}

func (updateIndexRules *UpdateIndexRules) AllPathsMustBeInIndex(paths []string) error {
	for _, path := range paths {
		index, err := updateIndexRules.indexRepository.Read()
		if err != nil {
			return err
		}
		if !index.HasEntry(path) {
			return errors.New(fmt.Sprintf("Path does not exist: %s", path))
		}
	}
	return nil
}

func (updateIndexRules *UpdateIndexRules) PathsMustNotDuplicate(paths []string) error {
	pathSet := make(map[string]bool)
	for _, path := range paths {
		if _, exists := pathSet[path]; exists {
			return errors.New(fmt.Sprintf("Duplicate path found: %s", path))
		}
		pathSet[path] = true
	}
	return nil
}
