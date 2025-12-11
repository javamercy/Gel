package vcs

import (
	"Gel/domain"
	"path"
)

type EntryProcessor = func(entry *domain.TreeEntry, relativePath string) error
type WalkOptions struct {
	Recursive    bool
	IncludeTrees bool
	OnlyTrees    bool
}

type TreeWalker struct {
	objectService *ObjectService
	options       WalkOptions
	processor     EntryProcessor
}

func NewTreeWalker(objectService *ObjectService, options WalkOptions, processor EntryProcessor) *TreeWalker {
	return &TreeWalker{
		objectService: objectService,
		options:       options,
		processor:     processor,
	}
}

func (treeWalker *TreeWalker) Walk(treeHash, prefix string) error {
	entries, err := treeWalker.objectService.ReadTreeAndDeserializeEntries(treeHash)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		relativePath := path.Join(prefix, entry.Name)

		isTree := entry.Mode.IsDirectory()
		shouldProcess := treeWalker.shouldProcessEntry(isTree)

		if shouldProcess {
			if err := treeWalker.processor(entry, relativePath); err != nil {
				return err
			}
		}

		if isTree && treeWalker.options.Recursive {
			if err := treeWalker.Walk(entry.Hash, relativePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func (treeWalker *TreeWalker) shouldProcessEntry(isTree bool) bool {
	if treeWalker.options.OnlyTrees {
		return isTree
	}

	if !isTree {
		return true
	}

	return treeWalker.options.IncludeTrees || !treeWalker.options.Recursive
}
