package gel

import (
	"Gel/domain"
	"Gel/internal/pathspec"
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
	pathResolver *pathspec.PathResolver,
) *AddService {
	return &AddService{
		indexService:       indexService,
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (a *AddService) Add(writer io.Writer, pathspecs []string, dryRun, verbose bool) error {
	index, err := a.indexService.Read()
	if err != nil {
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

	addedFiles, err := a.updateIndexService.UpdateIndex(pathsToAdd, true, false, !dryRun)
	if err != nil {
		return err
	}
	removedFiles, err := a.updateIndexService.UpdateIndex(pathsToRemove, false, true, !dryRun)
	if err != nil {
		return err
	}

	if verbose || dryRun {
		return addWithDryRun(writer, addedFiles, removedFiles)
	}
	return nil
}

func collectPaths(index *domain.Index, resolvedPaths []pathspec.ResolvedPath) (
	[]string, []string, error,
) {
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
			} else {
				prefix := resolved.NormalizedScope
				if prefix != "" {
					prefix += "/"
				}
				indexEntries = index.FindEntriesByPathPrefix(prefix)
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
			return nil, nil, fmt.Errorf(
				"pathspec '%s' did not match any files", resolved.NormalizedScope,
			)
		}

	}
	return pathsToAdd, pathsToRemove, nil
}

func addWithDryRun(writer io.Writer, pathsToAdd, pathsToRemove []string) error {
	for _, path := range pathsToAdd {
		if _, err := writer.Write([]byte(fmt.Sprintf("add '%s'\n", path))); err != nil {
			return err
		}
	}
	for _, path := range pathsToRemove {
		if _, err := writer.Write([]byte(fmt.Sprintf("remove '%s'\n", path))); err != nil {
			return err
		}
	}
	return nil
}
