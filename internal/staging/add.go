package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"fmt"
	"io"
)

type AddOptions struct {
	DryRun  bool
	Verbose bool
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

func (a *AddService) Add(writer io.Writer, pathspecs []string, options AddOptions) error {
	index, err := a.indexService.Read()
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	resolvedPaths, err := a.pathResolver.Resolve(pathspecs)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	pathsToAdd, pathsToRemove, err := a.collectPaths(index, resolvedPaths)
	if err != nil {
		return err
	}

	addedFiles, err := a.updateIndexService.UpdateIndex(
		pathsToAdd, UpdateIndexOptions{
			Add:    true,
			Remove: false,
			Write:  !options.DryRun,
		},
	)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	removedFiles, err := a.updateIndexService.UpdateIndex(
		pathsToRemove, UpdateIndexOptions{
			Add:    false,
			Remove: true,
			Write:  !options.DryRun,
		},
	)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	if options.Verbose || options.DryRun {
		return a.addWithDryRun(writer, addedFiles, removedFiles)
	}
	return nil
}

func (a *AddService) collectPaths(index *domain.Index, resolvedPaths []core.ResolvedPath) ([]string, []string, error) {
	var pathsToAdd []string
	var pathsToRemove []string

	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			pathsToAdd = append(pathsToAdd, path.String())
		}

		var indexEntries []*domain.IndexEntry

		switch resolved.Type {
		case core.File, core.NonExistent:
			if entry, _ := index.FindEntry(resolved.NormalizedScope); entry != nil {
				indexEntries = []*domain.IndexEntry{entry}
			} else {
				prefix := resolved.NormalizedScope
				if prefix != "" {
					prefix += "/"
				}
				indexEntries = index.FindEntriesByPathPrefix(prefix)
			}
		case core.Directory:
			prefix := resolved.NormalizedScope
			if prefix != "" {
				prefix += "/"
			}
			indexEntries = index.FindEntriesByPathPrefix(prefix)
		case core.GlobPattern:
			indexEntries = index.FindEntriesByPathPattern(resolved.NormalizedScope)
		}

		for _, entry := range indexEntries {
			if !resolved.NormalizedPaths[entry.Path] {
				pathsToRemove = append(pathsToRemove, entry.Path.String())
			}
		}

		if len(resolved.NormalizedPaths) == 0 && len(indexEntries) == 0 {
			return nil, nil, fmt.Errorf("'%s': %w", resolved.NormalizedScope, ErrPathDidNotMatch)
		}

	}
	return pathsToAdd, pathsToRemove, nil
}

func (a *AddService) addWithDryRun(writer io.Writer, pathsToAdd, pathsToRemove []string) error {
	for _, path := range pathsToAdd {
		if _, err := writer.Write([]byte(fmt.Sprintf("add '%s'\n", path))); err != nil {
			return fmt.Errorf("failed to write add message for '%s': %w", path, err)
		}
	}
	for _, path := range pathsToRemove {
		if _, err := writer.Write([]byte(fmt.Sprintf("remove '%s'\n", path))); err != nil {
			return fmt.Errorf("failed to write remove message for '%s': %w", path, err)
		}
	}
	return nil
}
