package gel

import (
	"Gel/domain"
	"Gel/internal/gel/validate"
	"time"
)

type ReadTreeService struct {
	indexService  *IndexService
	objectService *ObjectService
}

func NewReadTreeService(indexService *IndexService, objectService *ObjectService) *ReadTreeService {
	return &ReadTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (readTreeService *ReadTreeService) ReadTree(hash string) error {
	if err := validate.Hash(hash); err != nil {
		return err
	}

	var indexEntries []*domain.IndexEntry
	processor := func(entry domain.TreeEntry, relPath string) error {
		indexEntry, err := domain.NewIndexEntry(
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

		if err != nil {
			return err
		}
		indexEntries = append(indexEntries, indexEntry)
		return nil
	}

	options := WalkOptions{
		Recursive:    true,
		IncludeTrees: false,
		OnlyTrees:    false,
	}

	treeWalker := NewTreeWalker(readTreeService.objectService, options)
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
