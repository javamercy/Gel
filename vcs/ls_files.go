package vcs

import (
	"Gel/domain"
	"Gel/storage"
	"fmt"
	"io"
	"os"
)

type LsFilesService struct {
	indexService      *IndexService
	filesystemStorage *storage.FilesystemStorage
	objectService     *ObjectService
}

func NewLsFilesService(indexService *IndexService, filesystemStorage *storage.FilesystemStorage, objectService *ObjectService) *LsFilesService {
	return &LsFilesService{
		indexService:      indexService,
		filesystemStorage: filesystemStorage,
		objectService:     objectService,
	}
}

func (l *LsFilesService) LsFiles(writer io.Writer, cached, stage, modified, deleted bool) error {
	index, err := l.indexService.Read()
	if err != nil {
		return err
	}

	entries := index.Entries

	if stage {
		return l.LsFilesWithStage(writer, entries)
	} else if cached {
		return l.LsFilesWithCache(writer, entries)
	} else if modified {
		return l.LsFilesWithModified(writer, entries)
	} else if deleted {
		return l.LsFilesWithDeleted(writer, entries)
	}
	return l.LsFilesWithCache(writer, entries)
}

func (l *LsFilesService) LsFilesWithStage(writer io.Writer, entries []*domain.IndexEntry) error {
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

func (l *LsFilesService) LsFilesWithCache(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithModified(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := l.filesystemStorage.Exists(entry.Path)
		if !exists {
			continue
		}

		isModified := l.isModified(entry)
		if !isModified {
			continue
		}
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
			return err
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithDeleted(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		exists := l.filesystemStorage.Exists(entry.Path)
		if !exists {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *LsFilesService) isModified(entry *domain.IndexEntry) bool {
	path := entry.Path
	stat, err := os.Stat(path)

	if err != nil {
		return false
	}
	if uint32(stat.Size()) != entry.Size {
		return true
	}
	if !stat.ModTime().Equal(entry.UpdatedTime) {
		currentHash, err := l.objectService.ComputeHash(path)
		if err != nil {
			return false
		}
		return currentHash != entry.Hash
	}
	return false
}
