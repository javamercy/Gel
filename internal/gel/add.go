package gel

import (
	"Gel/domain"
	"Gel/internal/pathspec"
	"errors"
	"fmt"
	"io"
)

type AddService struct {
	indexService       *IndexService
	updateIndexService *UpdateIndexService
	pathResolver       *pathspec.PathResolver
}

func NewAddService(
	indexService *IndexService,
	updateIndexService *UpdateIndexService,
	pathResolver *pathspec.PathResolver) *AddService {
	return &AddService{
		indexService:       indexService,
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (a *AddService) Add(w io.Writer, pathspecs []string, dryRun bool) error {
	index, err := a.indexService.Read()
	if errors.Is(err, ErrIndexNotFound) {
		index = domain.NewEmptyIndex()
	} else if err != nil {
		return err
	}

	resolvedPaths, err := a.pathResolver.Resolve(pathspecs)
	if err != nil {
		return err
	}

	pathsToAdd, pathsToRemove, err := collectPaths(index, resolvedPaths)
	if err != nil {
		return err
	}

	if dryRun {
		return addWithDryRun(w, pathsToAdd, pathsToRemove)
	}

	if err := a.updateIndexService.UpdateIndex(pathsToAdd, true, false); err != nil {
		return err
	}
	if err := a.updateIndexService.UpdateIndex(pathsToRemove, false, true); err != nil {
		return err
	}
	return nil
}

func collectPaths(index *domain.Index, resolvedPaths []pathspec.ResolvedPath) ([]string, []string, error) {
	var pathsToAdd []string
	var pathsToRemove []string

	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			pathsToAdd = append(pathsToAdd, path)
		}

		var indexEntries []*domain.IndexEntry

		switch resolved.Type {
		case pathspec.File, pathspec.NonExistent:
			if entry, _ := index.FindEntry(resolved.NormalizedScope); entry != nil {
				indexEntries = []*domain.IndexEntry{entry}
			}
		case pathspec.Directory:
			prefix := resolved.NormalizedScope
			if prefix != "" {
				prefix += "/"
			}
			indexEntries = index.FindEntriesByPathPrefix(prefix)
		case pathspec.GlobPattern:
			indexEntries = index.FindEntriesByPathPattern(resolved.NormalizedScope)
		}

		for _, entry := range indexEntries {
			if !resolved.NormalizedPaths[entry.Path] {
				pathsToRemove = append(pathsToRemove, entry.Path)
			}
		}

		if len(resolved.NormalizedPaths) == 0 && len(indexEntries) == 0 {
			return nil, nil, fmt.Errorf("pathspec '%s' did not match any files", resolved.NormalizedScope)
		}

	}
	return pathsToAdd, pathsToRemove, nil
}

func addWithDryRun(w io.Writer, pathsToAdd, pathsToRemove []string) error {
	for _, path := range pathsToAdd {
		if _, err := w.Write([]byte(fmt.Sprintf("A  %s\n", path))); err != nil {
			return err
		}
	}
	for _, path := range pathsToRemove {
		if _, err := w.Write([]byte(fmt.Sprintf("D  %s\n", path))); err != nil {
			return err
		}
	}
	return nil
}
