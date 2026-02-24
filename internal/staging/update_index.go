package staging

import (
	"Gel/domain"
	"Gel/internal/core"
)

type UpdateIndexService struct {
	indexService      *core.IndexService
	objectService     *core.ObjectService
	hashObjectService *core.HashObjectService
	changeDetector    *core.ChangeDetector
}

func NewUpdateIndexService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	hashObjectService *core.HashObjectService,
	changeDetector *core.ChangeDetector,
) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		objectService:     objectService,
		hashObjectService: hashObjectService,
		changeDetector:    changeDetector,
	}
}

func (u *UpdateIndexService) UpdateIndex(paths []string, add, remove, write bool) (
	[]string, error,
) {
	index, err := u.indexService.Read()
	if err != nil {
		return nil, err
	}

	switch {
	case add:
		return u.updateIndexWithAdd(index, paths, write)
	case remove:
		return u.updateIndexWithRemove(index, paths, write)
	default:
		return nil, nil
	}
}

func (u *UpdateIndexService) updateIndexWithAdd(index *domain.Index, paths []string, write bool) (
	[]string, error,
) {
	var addedPaths []string
	for _, path := range paths {
		stat := domain.GetFileStatFromPath(path)
		entry, _ := index.FindEntry(path)

		var newEntry *domain.IndexEntry
		if entry != nil {
			changeResult, err := u.changeDetector.DetectFileChange(entry, stat)
			if err != nil {
				return nil, err
			}

			if !changeResult.IsModified {
				continue
			}

			addedPaths = append(addedPaths, path)

			if !write {
				continue
			}

			newEntry = domain.NewIndexEntry(
				path,
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
			hash, _, err := u.hashObjectService.HashObject(path, true)
			if err != nil {
				return nil, err
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
	return addedPaths, u.indexService.Write(index)
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
		return nil, err
	}
	return removedPaths, nil
}
