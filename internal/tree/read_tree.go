package tree

import (
	"Gel/domain"
	core2 "Gel/internal/core"
	"time"
)

type ReadTreeService struct {
	indexService  *core2.IndexService
	objectService *core2.ObjectService
}

func NewReadTreeService(indexService *core2.IndexService, objectService *core2.ObjectService) *ReadTreeService {
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

	options := core2.WalkOptions{
		Recursive:    true,
		IncludeTrees: false,
		OnlyTrees:    false,
	}

	treeWalker := core2.NewTreeWalker(readTreeService.objectService, options)
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
