package tree

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
	"time"
)

type ReadTreeService struct {
	indexService  *core.IndexService
	objectService *core.ObjectService
}

func NewReadTreeService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
) *ReadTreeService {
	return &ReadTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (readTreeService *ReadTreeService) ReadTree(hash domain2.Hash) error {
	var indexEntries []*domain2.IndexEntry
	processor := func(entry domain2.TreeEntry, path string) error {
		// TODO: fix here later
		normalizedPath, err := domain2.NewNormalizedPath("", path)
		if err != nil {
			return err
		}
		indexEntry := domain2.NewIndexEntry(
			normalizedPath,
			entry.Hash,
			0,
			entry.Mode.Uint32(),
			0,
			0,
			0,
			0,
			domain2.ComputeIndexFlags(normalizedPath.String(), 0),
			time.Time{},
			time.Time{},
		)
		indexEntries = append(indexEntries, indexEntry)
		return nil
	}

	options := core.WalkOptions{
		Recursive:    true,
		IncludeTrees: false,
		OnlyTrees:    false,
	}

	treeWalker := core.NewTreeWalker(readTreeService.objectService, options)
	err := treeWalker.Walk(hash, "", processor)
	if err != nil {
		return err
	}

	index := domain2.NewEmptyIndex()
	for _, entry := range indexEntries {
		index.AddEntry(entry)
	}
	return readTreeService.indexService.Write(index)
}
