package gel

import (
	"Gel/domain"
	"path"
)

type Processor = func(entry domain.TreeEntry, relPath string) error

type WalkOptions struct {
	Recursive    bool
	IncludeTrees bool
	OnlyTrees    bool
}

type TreeWalker struct {
	objectService *ObjectService
	options       WalkOptions
}

func NewTreeWalker(objectService *ObjectService, options WalkOptions) *TreeWalker {
	return &TreeWalker{
		objectService: objectService,
		options:       options,
	}
}

func (w *TreeWalker) Walk(hash, prefix string, processor Processor) error {
	entries, err := w.objectService.ReadTreeAndDeserializeEntries(hash)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		relPath := path.Join(prefix, entry.Name)
		isTree := entry.Mode.IsDirectory()
		shouldProcess := w.shouldProcess(isTree)

		if shouldProcess {
			if err := processor(entry, relPath); err != nil {
				return err
			}
		}
		if isTree && w.options.Recursive {
			if err := w.Walk(entry.Hash, relPath, processor); err != nil {
				return err
			}
		}
	}
	return nil
}

// shouldProcess determines if the given entry, based on its type, should be processed according to the walk options.
func (w *TreeWalker) shouldProcess(isTree bool) bool {
	if isTree {
		return w.options.IncludeTrees || w.options.OnlyTrees
	}
	return !w.options.OnlyTrees
}
