package tree

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
	"fmt"
	"io"
)

type LsTreeService struct {
	objectService *core.ObjectService
}

func NewLsTreeService(objectService *core.ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (l *LsTreeService) LsTree(writer io.Writer, hash domain2.Hash, recursive, showTrees, nameOnly bool) error {
	processor := func(entry domain2.TreeEntry, relPath string) error {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return err
		}
		if nameOnly {
			if _, err := fmt.Fprintln(writer, entry.Name); err != nil {
				return fmt.Errorf("failed to write entry name: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(
				writer,
				"%s %s %s\t%s\n",
				entry.Mode,
				objectType,
				entry.Hash,
				relPath,
			); err != nil {
				return fmt.Errorf("failed to write entry: %w", err)
			}
		}
		return nil
	}

	options := core.WalkOptions{
		Recursive:    recursive,
		IncludeTrees: showTrees,
		OnlyTrees:    false,
	}
	treeWalker := core.NewTreeWalker(l.objectService, options)
	return treeWalker.Walk(hash, "", processor)
}
