package vcs

import (
	"Gel/domain"
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

func (lsTreeService *LsTreeService) LsTree(writer io.Writer, treeHash string, recursive, showTrees bool) error {

	// TODO: validate hash string

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

	treeWalker := NewTreeWalker(lsTreeService.objectService, options, processor)
	return treeWalker.Walk(treeHash, "")
}
