package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"fmt"
)

type AddOptions struct {
	DryRun  bool
	Verbose bool
}

type AddResult struct {
	Added   []domain.AbsolutePath
	Removed []domain.AbsolutePath
	Error   error
}
type AddService struct {
	indexService       *core.IndexService
	updateIndexService *UpdateIndexService
	pathResolver       *core.PathResolver
}

func NewAddService(
	indexService *core.IndexService,
	updateIndexService *UpdateIndexService,
	pathResolver *core.PathResolver,
) *AddService {
	return &AddService{
		indexService:       indexService,
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
	}
}

func (a *AddService) Add(pathspecs []string, options AddOptions) AddResult {
	index, err := a.indexService.Read()
	if err != nil {
		return AddResult{Error: fmt.Errorf("add: %w", err)}
	}

	resolvedPaths, err := a.pathResolver.Resolve(pathspecs)
	if err != nil {
		return AddResult{Error: fmt.Errorf("add: %w", err)}
	}

	pathsToAdd, pathsToRemove, err := a.collectPaths(index, resolvedPaths)
	if err != nil {
		return AddResult{Error: fmt.Errorf("add: %w", err)}
	}

	if options.DryRun {
		return AddResult{Added: pathsToAdd, Removed: pathsToRemove}
	}

	addedFiles, err := a.updateIndexService.UpdateIndex(
		pathsToAdd, UpdateIndexOptions{
			Add:    true,
			Remove: false,
			Write:  !options.DryRun,
		},
	)
	if err != nil {
		return AddResult{Error: fmt.Errorf("add: %w", err)}
	}

	removedFiles, err := a.updateIndexService.UpdateIndex(
		pathsToRemove, UpdateIndexOptions{
			Add:    false,
			Remove: true,
			Write:  !options.DryRun,
		},
	)
	if err != nil {
		return AddResult{Error: fmt.Errorf("add: %w", err)}
	}
	if options.Verbose {
		return AddResult{Added: addedFiles, Removed: removedFiles}
	}
	return AddResult{}
}

func (a *AddService) collectPaths(index *domain.Index, resolvedPaths []core.ResolvedPath) (
	[]domain.AbsolutePath, []domain.AbsolutePath, error,
) {
	var pathsToAdd []domain.AbsolutePath
	var pathsToRemove []domain.AbsolutePath

	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			absolutePath, err := path.ToAbsolutePath()
			if err != nil {
				return nil, nil, fmt.Errorf("add: %w", err)
			}
			pathsToAdd = append(pathsToAdd, absolutePath)
		}

		var indexEntries []*domain.IndexEntry

		switch resolved.Type {
		case core.PathspecTypeFile, core.PathspecTypeNonExistent:
			if entry, _ := index.FindEntry(resolved.NormalizedScope); entry != nil {
				indexEntries = []*domain.IndexEntry{entry}
			} else {
				prefix := resolved.NormalizedScope
				if prefix != "" {
					prefix += "/"
				}
				indexEntries = index.FindEntriesByPathPrefix(prefix)
			}
		case core.PathspecTypeDirectory:
			prefix := resolved.NormalizedScope
			if prefix != "" {
				prefix += "/"
			}
			indexEntries = index.FindEntriesByPathPrefix(prefix)
		case core.PathspecTypeGlobPattern:
			indexEntries = index.FindEntriesByPathPattern(resolved.NormalizedScope)
		}

		for _, entry := range indexEntries {
			if !resolved.NormalizedPaths[entry.Path] {
				absolutePath, err := entry.Path.ToAbsolutePath()
				if err != nil {
					return nil, nil, fmt.Errorf("add: %w", err)
				}
				pathsToRemove = append(pathsToRemove, absolutePath)
			}
		}

		if len(resolved.NormalizedPaths) == 0 && len(indexEntries) == 0 {
			return nil, nil, fmt.Errorf("'%s': %w", resolved.NormalizedScope, ErrPathDidNotMatch)
		}

	}
	return pathsToAdd, pathsToRemove, nil
}
