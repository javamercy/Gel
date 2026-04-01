package staging

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
	"Gel/internal/validate"
	"errors"
	"fmt"
)

// UpdateIndexOptions controls update-index behavior.
type UpdateIndexOptions struct {
	// Add updates or inserts entries for the provided paths.
	Add bool
	// Remove deletes entries for the provided paths from the index.
	Remove bool
	// Write persists the updated index to disk when true.
	Write bool
}

// UpdateIndexService updates index entries from working tree files.
type UpdateIndexService struct {
	indexService      *core.IndexService
	objectService     *core.ObjectService
	hashObjectService *core.HashObjectService
	changeDetector    *core.ChangeDetector
	workspace         *domain2.Workspace
}

// NewUpdateIndexService creates an index update service with its dependencies.
func NewUpdateIndexService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	hashObjectService *core.HashObjectService,
	changeDetector *core.ChangeDetector,
	workspace *domain2.Workspace,
) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		objectService:     objectService,
		hashObjectService: hashObjectService,
		changeDetector:    changeDetector,
		workspace:         workspace,
	}
}

// UpdateIndex applies add/remove operations for normalized repository paths.
//
// At least one of options.Add or options.Remove must be enabled. Returned paths
// are the paths that were actually affected by the selected operation.
func (u *UpdateIndexService) UpdateIndex(
	paths []domain2.NormalizedPath,
	options UpdateIndexOptions,
) ([]domain2.NormalizedPath, error) {
	if !options.Add && !options.Remove {
		return nil, errors.New("update-index: must specify --add or --remove")
	}

	index, err := u.indexService.Read()
	if err != nil {
		return nil, fmt.Errorf("update-index: %w", err)
	}

	switch {
	case options.Add:
		return u.updateIndexWithAdd(index, paths, options.Write)
	case options.Remove:
		return u.updateIndexWithRemove(index, paths, options.Write)
	default:
		return nil, nil
	}
}

// updateIndexWithAdd stages file content for the given normalized paths.
// It computes object hashes, writes blob objects when requested, and updates
// entry metadata only for paths that are newly added or modified.
func (u *UpdateIndexService) updateIndexWithAdd(
	index *domain2.Index,
	paths []domain2.NormalizedPath,
	write bool,
) (
	[]domain2.NormalizedPath, error,
) {
	var addedPaths []domain2.NormalizedPath
	for _, path := range paths {
		absolutePath, err := path.ToAbsolutePath(u.workspace.RepoDir)
		if err != nil {
			return nil, fmt.Errorf("update-index: %w", err)
		}
		if err := validate.PathMustBeFile(absolutePath.String()); err != nil {
			return nil, fmt.Errorf("update-index: %w", err)
		}

		var newEntry *domain2.IndexEntry
		stat := domain2.GetFileStatFromPath(absolutePath)
		entry, _ := index.FindEntry(path)
		if entry != nil {
			changeResult, err := u.changeDetector.DetectFileChange(entry, stat)
			if err != nil {
				return nil, fmt.Errorf("update-index: %w", err)
			}

			if !changeResult.IsModified {
				continue
			}

			addedPaths = append(addedPaths, path)
			if !write {
				continue
			}
			if _, err := u.hashObjectService.HashObject(
				absolutePath, core.HashObjectOptions{Write: true},
			); err != nil {
				return nil, fmt.Errorf("update-index: %w", err)
			}

			index.RemoveEntry(path)
			newEntry = domain2.NewIndexEntry(
				path,
				changeResult.NewHash,
				stat.Size,
				domain2.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain2.ComputeIndexFlags(path.String(), 0),
				stat.CreatedTime,
				stat.UpdatedTime,
			)
		} else {
			hash, _, err := u.objectService.ComputeObjectHash(absolutePath)
			if err != nil {
				return nil, fmt.Errorf("update-index: %w", err)
			}

			addedPaths = append(addedPaths, path)
			if !write {
				continue
			}
			if _, err := u.hashObjectService.HashObject(
				absolutePath, core.HashObjectOptions{Write: true},
			); err != nil {
				return nil, fmt.Errorf("update-index: %w", err)
			}

			newEntry = domain2.NewIndexEntry(
				path,
				hash,
				stat.Size,
				domain2.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain2.ComputeIndexFlags(path.String(), 0),
				stat.CreatedTime,
				stat.UpdatedTime,
			)
		}
		index.SetEntry(newEntry)
	}
	if !write {
		return addedPaths, nil
	}

	err := u.indexService.Write(index)
	if err != nil {
		return nil, fmt.Errorf("update-index: %w", err)
	}
	return addedPaths, nil
}

// updateIndexWithRemove removes the given paths from the index.
// Missing paths are treated as no-op removals.
func (u *UpdateIndexService) updateIndexWithRemove(index *domain2.Index, paths []domain2.NormalizedPath, write bool) (
	[]domain2.NormalizedPath, error,
) {
	var removedPaths []domain2.NormalizedPath
	for _, path := range paths {
		if index.HasEntry(path) {
			removedPaths = append(removedPaths, path)
		}
		index.RemoveEntry(path)
	}
	if !write {
		return removedPaths, nil
	}
	if err := u.indexService.Write(index); err != nil {
		return nil, fmt.Errorf("update-index: %w", err)
	}
	return removedPaths, nil
}
