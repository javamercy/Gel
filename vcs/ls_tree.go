package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"strings"
)

type LsTreeService struct {
	objectService *ObjectService
}

func NewLsTreeService(objectService *ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (lsTreeService *LsTreeService) LsTree(treeHash string, recursive, showTrees bool) (string, error) {

	// TODO: validate hash string

	var result strings.Builder

	processor := func(entry domain.TreeEntry, relativePath string) error {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return err
		}

		result.WriteString(entry.Mode.String())
		result.WriteString(constant.SpaceStr)
		result.WriteString(string(objectType))
		result.WriteString(constant.SpaceStr)
		result.WriteString(entry.Hash)
		result.WriteString(constant.TabStr)
		result.WriteString(relativePath)
		result.WriteString(constant.NewLineStr)

		return nil
	}

	options := WalkOptions{
		Recursive:    recursive,
		IncludeTrees: showTrees,
		OnlyTrees:    false,
	}

	treeWalker := NewTreeWalker(lsTreeService.objectService, options, processor)
	err := treeWalker.Walk(treeHash, "")
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
