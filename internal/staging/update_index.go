package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/validate"
	"Gel/internal/workspace"
	"errors"
	"fmt"
)

type UpdateIndexOptions struct {
	Add    bool
	Remove bool
	Write  bool
}
type UpdateIndexService struct {
	indexService      *core.IndexService
	objectService     *core.ObjectService
	hashObjectService *core.HashObjectService
	changeDetector    *core.ChangeDetector
	workspaceProvider *workspace.Provider
}

func NewUpdateIndexService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	hashObjectService *core.HashObjectService,
	changeDetector *core.ChangeDetector,
	workspaceProvider *workspace.Provider,
) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		objectService:     objectService,
		hashObjectService: hashObjectService,
		changeDetector:    changeDetector,
		workspaceProvider: workspaceProvider,
	}
}

func (u *UpdateIndexService) UpdateIndex(paths []string, options UpdateIndexOptions) ([]string, error) {
	if !options.Add && !options.Remove {
		return nil, errors.New("update-index: must specify --add or --remove")
	}

	index, err := u.indexService.Read()
	if err != nil {
		return nil, fmt.Errorf("update index: %w", err)
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

func (u *UpdateIndexService) updateIndexWithAdd(index *domain.Index, paths []string, write bool) ([]string, error) {
	var addedPaths []string
	for _, path := range paths {
		if err := validate.PathMustBeFile(path); err != nil {
			return nil, fmt.Errorf("update index: %w", err)
		}

		var newEntry *domain.IndexEntry
		absPath, err := domain.NewAbsolutePath(path)
		if err != nil {
			return nil, fmt.Errorf("update-index: %w", err)
		}
		stat := domain.GetFileStatFromPath(absPath)
		entry, _ := index.FindEntry(path)
		if entry != nil {
			changeResult, err := u.changeDetector.DetectFileChange(entry, stat)
			if err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}

			if !changeResult.IsModified {
				continue
			}

			addedPaths = append(addedPaths, path)

			if !write {
				continue
			}

			if _, err := u.hashObjectService.HashObject(path, core.HashObjectOptions{Write: true}); err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}

			normalizedPath, err := absPath.ToNormalizedPath(u.workspaceProvider.GetWorkspace().RepoDir)
			if err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}
			index.RemoveEntry(normalizedPath.String())
			newEntry = domain.NewIndexEntry(
				normalizedPath,
				changeResult.NewHash,
				stat.Size,
				domain.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain.ComputeIndexFlags(path, 0),
				stat.CreatedTime,
				stat.UpdatedTime,
			)
		} else {
			absPath, err := domain.NewAbsolutePath(path)
			if err != nil {
				return nil, fmt.Errorf("update-index: %w", err)
			}
			hash, _, err := u.objectService.ComputeObjectHash(absPath)
			if err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}

			addedPaths = append(addedPaths, path)

			if !write {
				continue
			}

			if _, err := u.hashObjectService.HashObject(path, core.HashObjectOptions{Write: true}); err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}

			normalizedPath, err := absPath.ToNormalizedPath(u.workspaceProvider.GetWorkspace().RepoDir)
			if err != nil {
				return nil, fmt.Errorf("update index: %w", err)
			}
			newEntry = domain.NewIndexEntry(
				normalizedPath,
				hash,
				stat.Size,
				domain.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain.ComputeIndexFlags(path, 0),
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
		return nil, fmt.Errorf("update index: %w", err)
	}
	return addedPaths, nil
}

func (u *UpdateIndexService) updateIndexWithRemove(index *domain.Index, paths []string, write bool) (
	[]string, error,
) {
	var removedPaths []string
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
		return nil, fmt.Errorf("update index: %w", err)
	}
	return removedPaths, nil
}
