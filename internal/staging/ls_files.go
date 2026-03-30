package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"errors"
	"fmt"
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

func (l *LsFilesService) LsFiles(pathspec string, options LsFilesOptions) ([]string, error) {
	if !options.Stage && !options.Cached && !options.Modified && !options.Deleted {
		return nil, fmt.Errorf("ls-files: must specify --stage, --cached, --modified, or --deleted")
	}

	index, err := l.indexService.Read()
	if err != nil {
		return nil, fmt.Errorf("ls-files: %w", err)
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
		return l.LsFilesWithStage(entries), nil
	case options.Cached:
		return l.LsFilesWithCached(entries), nil
	case options.Modified:
		return l.LsFilesWithModified(entries)
	case options.Deleted:
		return l.LsFilesWithDeleted(entries)
	}
	return nil, nil
}

func (l *LsFilesService) LsFilesWithStage(entries []*domain.IndexEntry) []string {
	files := make([]string, len(entries))
	for i, entry := range entries {
		files[i] = fmt.Sprintf(
			"%s %s %d\t%s",
			domain.ParseFileMode(entry.Mode),
			entry.Hash,
			entry.GetStage(),
			entry.Path,
		)
	}
	return files
}

func (l *LsFilesService) LsFilesWithCached(entries []*domain.IndexEntry) []string {
	files := make([]string, len(entries))
	for i, entry := range entries {
		files[i] = entry.Path.String()
	}
	return files
}

func (l *LsFilesService) LsFilesWithModified(entries []*domain.IndexEntry) ([]string, error) {
	files := make([]string, 0)
	for _, entry := range entries {
		absolutePath, err := entry.Path.ToAbsolutePath()
		if err != nil {
			return nil, fmt.Errorf("ls-files: %w", err)
		}
		stat := domain.GetFileStatFromPath(absolutePath)
		changeResult, err := l.changeDetector.DetectFileChange(entry, stat)
		if err != nil {
			return nil, err
		}
		if changeResult.IsModified {
			files = append(files, entry.Path.String())
		}
	}
	return files, nil
}

func (l *LsFilesService) LsFilesWithDeleted(entries []*domain.IndexEntry) ([]string, error) {
	files := make([]string, 0)
	for _, entry := range entries {
		absolutePath, err := entry.Path.ToAbsolutePath()
		if err != nil {
			return nil, fmt.Errorf("ls-files: %w", err)
		}
		_, err = os.Stat(absolutePath.String())
		switch {
		case errors.Is(err, os.ErrNotExist):
			files = append(files, entry.Path.String())
		case err != nil:
			return nil, fmt.Errorf("ls-files: failed to stat file '%s': %w", entry.Path, err)
		}
	}
	return files, nil
}
