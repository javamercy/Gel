package rules

import (
	"Gel/src/gel/persistence/repositories"
	"errors"
	"fmt"
)

type CatFileRules struct {
	objectRepository repositories.IObjectRepository
}

func NewCatFileRules(objectRepository repositories.IObjectRepository) *CatFileRules {
	return &CatFileRules{
		objectRepository: objectRepository,
	}
}

func (rules *CatFileRules) ObjectMustExist(hash string) error {
	if !rules.objectRepository.Exists(hash) {
		return errors.New(fmt.Sprintf("fatal: Not a valid object name %s", hash))
	}
	return nil
}
