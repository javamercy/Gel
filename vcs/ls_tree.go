package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
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

		if _, err := io.WriteString(writer, entry.Mode.String()); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, string(objectType)); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, entry.Hash); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, constant.TabStr); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, relativePath); err != nil {
			return err
		}
		if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
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
