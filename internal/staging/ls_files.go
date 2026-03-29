package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	globPatterns string = "*?[]"
)

type LsFilesOptions struct {
	Cached   bool
	Stage    bool
	Modified bool
	Deleted  bool
}
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

func (l *LsFilesService) LsFiles(writer io.Writer, pathspec string, options LsFilesOptions) error {
	if !options.Stage && !options.Cached && !options.Modified && !options.Deleted {
		return fmt.Errorf("ls-files: must specify --stage, --cached, --modified, or --deleted")
	}

	index, err := l.indexService.Read()
	if err != nil {
		return fmt.Errorf("ls-files: %w", err)
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

	switch {
	case options.Stage:
		return l.LsFilesWithStage(writer, entries)
	case options.Cached:
		return l.LsFilesWithCached(writer, entries)
	case options.Modified:
		return l.LsFilesWithModified(writer, entries)
	case options.Deleted:
		return l.LsFilesWithDeleted(writer, entries)
	}
	return nil
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
			return fmt.Errorf("ls-files: failed to write entry: %w", err)
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithCached(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
			return fmt.Errorf("ls-files: failed to write entry: %w", err)
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithModified(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		stat := domain.GetFileStatFromPath(entry.Path.ToAbsolutePath())
		changeResult, err := l.changeDetector.DetectFileChange(entry, stat)
		if err != nil {
			return err
		}
		if changeResult.IsModified {
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return fmt.Errorf("ls-files: failed to write entry: %w", err)
			}
		}
	}
	return nil
}

func (l *LsFilesService) LsFilesWithDeleted(writer io.Writer, entries []*domain.IndexEntry) error {
	for _, entry := range entries {
		_, err := os.Stat(entry.Path.ToAbsolutePath().String())
		switch {
		case errors.Is(err, os.ErrNotExist):
			if _, err := fmt.Fprintf(writer, "%s\n", entry.Path); err != nil {
				return fmt.Errorf("ls-files: failed to write entry: %w", err)
			}
		case err != nil:
			return fmt.Errorf("ls-files: failed to stat file '%s': %w", entry.Path, err)
		}
	}
	return nil
}
