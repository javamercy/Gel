package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	globPatterns string = "*?[]"
)

type LsFilesService struct {
	indexService      *core.IndexService
	objectService     *core.ObjectService
	hashObjectService *core.HashObjectService
}

func NewLsFilesService(indexService *core.IndexService, objectService *core.ObjectService) *LsFilesService {
	return &LsFilesService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (l *LsFilesService) LsFiles(
	writer io.Writer, pathspec string, cached, stage, modified, deleted bool,
) error {
	index, err := l.indexService.Read()
	if err != nil {
		return err
	}

	var entries []*domain.IndexEntry
	if pathspec != "" {
		if strings.ContainsAny(pathspec, globPatterns) {
			entries = index.FindEntriesByPathPattern(pathspec)
		} else {
			entries = index.FindEntriesByPathPrefix(pathspec)
		}
	} else {
		entries = index.Entries
	}

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
		if _, err := fmt.Fprintf(
			writer,
			"%s %s %d\t%s\n",
			domain.ParseFileMode(entry.Mode),
			entry.Hash,
			entry.GetStage(),
			entry.Path,
		); err != nil {
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
		fileInfo, err := os.Stat(entry.Path)
		if err != nil {
			continue
		}
		if uint32(fileInfo.Size()) != entry.Size {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return err
			}
			return nil
		}
		if !fileInfo.ModTime().Equal(entry.UpdatedTime) {
			currentHash, _, err := l.hashObjectService.HashObject(entry.Path, false)
			if err != nil {
				return err
			}
			if currentHash != entry.Hash {
				if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithDeleted(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		_, err := os.Stat(entry.Path)
		if err != nil {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return err
			}
		}
	}
	return nil
}
