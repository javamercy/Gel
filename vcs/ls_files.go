package vcs

import (
	"Gel/domain"
	"fmt"
	"io"
	"os"
)

type LsFilesService struct {
	indexService      *IndexService
	filesystemService *FilesystemService
	objectService     *ObjectService
}

func NewLsFilesService(indexService *IndexService, filesystemService *FilesystemService, objectService *ObjectService) *LsFilesService {
	return &LsFilesService{
		indexService:      indexService,
		filesystemService: filesystemService,
		objectService:     objectService,
	}
}

func (lsFilesService *LsFilesService) LsFiles(writer io.Writer, cached, stage, modified, deleted bool) error {
	index, err := lsFilesService.indexService.Read()
	if err != nil {
		return err
	}

	entries := index.Entries

	if stage {
		return lsFilesService.LsFilesWithStage(writer, entries)
	} else if cached {
		return lsFilesService.LsFilesWithCache(writer, entries)
	} else if modified {
		return lsFilesService.LsFilesWithModified(writer, entries)
	} else if deleted {
		return lsFilesService.LsFilesWithDeleted(writer, entries)
	}
	return lsFilesService.LsFilesWithCache(writer, entries)
}

func (lsFilesService *LsFilesService) LsFilesWithStage(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer,
			"%s %s %d\t%s\n",
			domain.ParseFileMode(entry.Mode),
			entry.Hash,
			entry.GetStage(),
			entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithCache(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithModified(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			continue
		}

		isModified := lsFilesService.isModified(entry)
		if !isModified {
			continue
		}
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) LsFilesWithDeleted(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := lsFilesService.filesystemService.Exists(entry.Path)
		if !exists {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (lsFilesService *LsFilesService) isModified(entry *domain.IndexEntry) bool {
	path := entry.Path
	stat, err := os.Stat(path)

	if err != nil {
		return false
	}
	if uint32(stat.Size()) != entry.Size {
		return true
	}
	if !stat.ModTime().Equal(entry.UpdatedTime) {
		currentHash, err := lsFilesService.objectService.HashObject(path)
		if err != nil {
			return false
		}
		return currentHash != entry.Hash
	}
	return false
}
