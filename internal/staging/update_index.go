package staging

import (
	"Gel/domain"
	core2 "Gel/internal/core"
)

type UpdateIndexService struct {
	indexService      *core2.IndexService
	hashObjectService *core2.HashObjectService
	objectService     *core2.ObjectService
}

func NewUpdateIndexService(
	indexService *core2.IndexService, hashObjectService *core2.HashObjectService, objectService *core2.ObjectService,
) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		hashObjectService: hashObjectService,
		objectService:     objectService,
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
		fileStatInfo := domain.GetFileStatFromPath(path)

		if entry, _ := index.FindEntry(path); entry != nil {
			if entry.UpdatedTime.Equal(fileStatInfo.UpdatedTime) && entry.Size == fileStatInfo.Size {
				continue
			}
		}
		if !write {
			addedPaths = append(addedPaths, path)
			continue
		}

		hash, _, err := u.hashObjectService.HashObject(path, write)
		if err != nil {
			return nil, err
		}

		size, err := u.objectService.GetObjectSize(hash)
		if err != nil {
			return nil, err
		}

		newEntry := domain.NewIndexEntry(
			path,
			hash,
			size,
			domain.ParseFileModeFromOsMode(fileStatInfo.Mode).Uint32(),
			fileStatInfo.Device,
			fileStatInfo.Inode,
			fileStatInfo.UserId,
			fileStatInfo.GroupId,
			domain.ComputeIndexFlags(path, 0),
			fileStatInfo.CreatedTime,
			fileStatInfo.UpdatedTime,
		)
		addedPaths = append(addedPaths, path)
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
