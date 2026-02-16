package tree

import (
	"Gel/domain"
	"Gel/internal/gel/core"
	"time"
)

type ReadTreeService struct {
	indexService  *core.IndexService
	objectService *core.ObjectService
}

func NewReadTreeService(indexService *core.IndexService, objectService *core.ObjectService) *ReadTreeService {
	return &ReadTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (readTreeService *ReadTreeService) ReadTree(hash string) error {
	var indexEntries []*domain.IndexEntry
	processor := func(entry domain.TreeEntry, relPath string) error {
		indexEntry := domain.NewIndexEntry(
			relPath,
			entry.Hash,
			0,
			entry.Mode.Uint32(),
			0,
			0,
			0,
			0,
			domain.ComputeIndexFlags(relPath, 0),
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

	index := domain.NewEmptyIndex()
	for _, entry := range indexEntries {
		index.AddEntry(entry)
	}
	return readTreeService.indexService.Write(index)
}
