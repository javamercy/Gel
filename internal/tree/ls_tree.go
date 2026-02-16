package tree

import (
	"Gel/domain"
	core2 "Gel/internal/core"
	"fmt"
	"io"
)

type LsTreeService struct {
	objectService *core2.ObjectService
}

func NewLsTreeService(objectService *core2.ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (l *LsTreeService) LsTree(writer io.Writer, hash string, recursive, showTrees, nameOnly bool) error {
	processor := func(entry domain.TreeEntry, relPath string) error {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return err
		}
		if nameOnly {
			if _, err := fmt.Fprintln(writer, entry.Name); err != nil {
				return err
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
				return err
			}
		}
		return nil
	}

	options := core2.WalkOptions{
		Recursive:    recursive,
		IncludeTrees: showTrees,
		OnlyTrees:    false,
	}
	treeWalker := core2.NewTreeWalker(l.objectService, options)
	return treeWalker.Walk(hash, "", processor)
}
