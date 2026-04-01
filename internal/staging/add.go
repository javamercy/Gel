package staging

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
	"fmt"
)

type AddOptions struct {
	DryRun  bool
	Verbose bool
}

type AddResult struct {
	Added   []domain2.NormalizedPath
	Removed []domain2.NormalizedPath
	Error   error
}
type AddService struct {
	indexService       *core.IndexService
	updateIndexService *UpdateIndexService
	pathResolver       *core.PathResolver
	workspace          *domain2.Workspace
}

func NewAddService(
	indexService *core.IndexService,
	updateIndexService *UpdateIndexService,
	pathResolver *core.PathResolver,
	workspace *domain2.Workspace,
) *AddService {
	return &AddService{
		indexService:       indexService,
		updateIndexService: updateIndexService,
		pathResolver:       pathResolver,
		workspace:          workspace,
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

func (a *AddService) collectPaths(
	index *domain2.Index,
	resolvedPaths []core.ResolvedPath,
) (
	[]domain2.NormalizedPath, []domain2.NormalizedPath, error,
) {
	var pathsToAdd []domain2.NormalizedPath
	var pathsToRemove []domain2.NormalizedPath

	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			pathsToAdd = append(pathsToAdd, path)
		}

		var indexEntries []*domain2.IndexEntry
		switch resolved.Type {
		case core.PathspecTypeFile, core.PathspecTypeNonExistent:
			normalizedScope, err := domain2.NewNormalizedPathUnchecked(resolved.NormalizedScope)
			if err != nil {
				return nil, nil, fmt.Errorf("add: %w", err)
			}
			if entry, _ := index.FindEntry(normalizedScope); entry != nil {
				indexEntries = []*domain2.IndexEntry{entry}
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
				pathsToRemove = append(pathsToRemove, entry.Path)
			}
		}
		if len(resolved.NormalizedPaths) == 0 && len(indexEntries) == 0 {
			return nil, nil, fmt.Errorf("'%s': %w", resolved.NormalizedScope, ErrPathDidNotMatch)
		}

	}
	return pathsToAdd, pathsToRemove, nil
}
