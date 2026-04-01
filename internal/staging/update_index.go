package staging

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/validate"
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
	workspace         *domain.Workspace
}

func NewUpdateIndexService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	hashObjectService *core.HashObjectService,
	changeDetector *core.ChangeDetector,
	workspace *domain.Workspace,
) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		objectService:     objectService,
		hashObjectService: hashObjectService,
		changeDetector:    changeDetector,
		workspace:         workspace,
	}
}

func (u *UpdateIndexService) UpdateIndex(
	paths []domain.NormalizedPath,
	options UpdateIndexOptions,
) ([]domain.NormalizedPath, error) {
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

func (u *UpdateIndexService) updateIndexWithAdd(
	index *domain.Index,
	paths []domain.NormalizedPath,
	write bool,
) (
	[]domain.NormalizedPath, error,
) {
	var addedPaths []domain.NormalizedPath
	for _, path := range paths {
		if err := validate.PathMustBeFile(path.String()); err != nil {
			return nil, fmt.Errorf("update-index: %w", err)
		}

		var newEntry *domain.IndexEntry
		absolutePath, err := path.ToAbsolutePath(u.workspace.RepoDir)
		if err != nil {
			return nil, fmt.Errorf("update-index: %w", err)
		}

		stat := domain.GetFileStatFromPath(absolutePath)
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
			newEntry = domain.NewIndexEntry(
				path,
				changeResult.NewHash,
				stat.Size,
				domain.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain.ComputeIndexFlags(path.String(), 0),
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

			newEntry = domain.NewIndexEntry(
				path,
				hash,
				stat.Size,
				domain.ParseFileModeFromOsMode(stat.Mode).Uint32(),
				stat.Device,
				stat.Inode,
				stat.UserId,
				stat.GroupId,
				domain.ComputeIndexFlags(path.String(), 0),
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

func (u *UpdateIndexService) updateIndexWithRemove(index *domain.Index, paths []domain.NormalizedPath, write bool) (
	[]domain.NormalizedPath, error,
) {
	var removedPaths []domain.NormalizedPath
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
