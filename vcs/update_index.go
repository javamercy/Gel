package vcs

import (
	"Gel/core/util"
	"Gel/domain"
	"errors"
)

type UpdateIndexService struct {
	indexService      *IndexService
	hashObjectService *HashObjectService
	objectService     *ObjectService
}

func NewUpdateIndexService(indexService *IndexService, hashObjectService *HashObjectService, objectService *ObjectService) *UpdateIndexService {
	return &UpdateIndexService{
		indexService:      indexService,
		hashObjectService: hashObjectService,
		objectService:     objectService,
	}
}

func (u *UpdateIndexService) UpdateIndex(paths []string, add, remove bool) error {
	index, err := u.indexService.Read()
	if errors.Is(err, ErrIndexNotFound) {
		index = domain.NewEmptyIndex()
	} else if err != nil {
		return err
	}

	switch {
	case add:
		return u.updateIndexWithAdd(index, paths)
	case remove:
		return u.updateIndexWithRemove(index, paths)
	default:
		return nil
	}
}

func (u *UpdateIndexService) updateIndexWithAdd(index *domain.Index, paths []string) error {
	for _, path := range paths {
		fileStatInfo := util.GetFileStatFromPath(path)
		hash, _, err := u.hashObjectService.HashObject(path, true)
		if err != nil {
			return err
		}

		size, err := u.objectService.GetObjectSize(hash)
		if err != nil {
			return err
		}

		newEntry, err := domain.NewIndexEntry(
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
			fileStatInfo.UpdatedTime)

		if err != nil {
			return err
		}

		index.SetEntry(newEntry)
	}
	return u.indexService.Write(index)
}

func (u *UpdateIndexService) updateIndexWithRemove(index *domain.Index, paths []string) error {
	for _, path := range paths {
		index.RemoveEntry(path)
	}

	err := u.indexService.Write(index)
	if err != nil {
		return err
	}

	return nil
}
