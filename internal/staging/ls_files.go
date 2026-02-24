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
	indexService   *core.IndexService
	objectService  *core.ObjectService
	changeDetector *core.ChangeDetector
}

func NewLsFilesService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	changeDetector *core.ChangeDetector,
) *LsFilesService {
	return &LsFilesService{
		indexService:   indexService,
		objectService:  objectService,
		changeDetector: changeDetector,
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
		stat := domain.GetFileStatFromPath(entry.Path)
		changeResult, err := l.changeDetector.DetectFileChange(entry, stat)
		if err != nil {
			return err
		}
		if changeResult.IsModified {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return err
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
