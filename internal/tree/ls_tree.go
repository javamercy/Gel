package tree

import (
	"Gel/internal/core"
	"Gel/internal/domain"
	"fmt"
)

type LsTreeOptions struct {
	Recursive bool
	ShowTrees bool
	NameOnly  bool
}

type LsTreeService struct {
	objectService *core.ObjectService
}

func NewLsTreeService(objectService *core.ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (l *LsTreeService) LsTree(hash domain.Hash, options LsTreeOptions) ([]string, error) {
	var contents []string
	processor := func(entry domain.TreeEntry, relPath string) error {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return err
		}
		if options.NameOnly {
			contents = append(contents, relPath)
		} else {
			result := fmt.Sprintf(
				"%s %s %s\t%s",
				entry.Mode, objectType, entry.Hash, relPath,
			)
			contents = append(contents, result)
		}
		return nil
	}

	walkOptions := core.WalkOptions{
		Recursive:    options.Recursive,
		IncludeTrees: options.ShowTrees,
		OnlyTrees:    false,
	}
	treeWalker := core.NewTreeWalker(l.objectService, walkOptions)
	if err := treeWalker.Walk(hash, "", processor); err != nil {
		return nil, err
	}
	return contents, nil
}
