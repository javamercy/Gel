package vcs

import (
	"Gel/domain"
	"Gel/vcs/validate"
	"fmt"
	"io"
)

type LsTreeService struct {
	objectService *ObjectService
}

func NewLsTreeService(objectService *ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (lsTreeService *LsTreeService) LsTree(writer io.Writer, hash string, recursive, showTrees bool) error {
	if err := validate.Hash(hash); err != nil {
		return err
	}

	processor := func(entry domain.TreeEntry, relativePath string) error {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer,
			"%s %s %s\t%s\n",
			entry.Mode,
			objectType,
			entry.Hash,
			relativePath); err != nil {
			return err
		}
		return nil
	}

	options := WalkOptions{
		Recursive:    recursive,
		IncludeTrees: showTrees,
		OnlyTrees:    false,
	}

	treeWalker := NewTreeWalker(lsTreeService.objectService, options)
	return treeWalker.Walk(hash, "", processor)
}
