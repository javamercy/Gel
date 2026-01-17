package vcs

import (
	"Gel/core/util"
	"Gel/domain"
	"Gel/vcs/validate"
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

	processor := func(entry domain.TreeEntry, relativePath string) error {
		fileStatInfo := util.GetFileStatFromPath(relativePath)

		indexEntry, err := domain.NewIndexEntry(
			relativePath,
			entry.Hash,
			fileStatInfo.Size,
			entry.Mode.Uint32(),
			fileStatInfo.Device,
			fileStatInfo.Inode,
			fileStatInfo.UserId,
			fileStatInfo.GroupId,
			domain.ComputeIndexFlags(relativePath, 0),
			time.Now(),
			time.Now())

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

	treeWalker := NewTreeWalker(readTreeService.objectService, options, processor)
	err := treeWalker.Walk(hash, "")

	if err != nil {
		return err
	}

	index := domain.NewEmptyIndex()
	for _, entry := range indexEntries {
		index.AddEntry(entry)
	}

	return readTreeService.indexService.Write(index)
}
