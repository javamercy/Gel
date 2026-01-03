package vcs

import (
	"Gel/core/encoding"
	"Gel/core/util"
	"Gel/domain"
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

func (updateIndexService *UpdateIndexService) UpdateIndex(paths []string, add, remove bool) error {
	index, err := updateIndexService.indexService.Read()
	if err != nil {
		index = domain.NewEmptyIndex()
	}

	if add {
		return updateIndexService.updateIndexWithAdd(index, paths)
	} else if remove {
		return updateIndexService.updateIndexWithRemove(index, paths)
	}
	return nil
}

func (updateIndexService *UpdateIndexService) updateIndexWithAdd(index *domain.Index, paths []string) error {

	hashMap, err := updateIndexService.hashObjectService.HashObject(paths, true)
	if err != nil {
		return err
	}

	for _, p := range paths {

		fileStatInfo := util.GetFileStatFromPath(p)

		blobHash := hashMap[p]
		size, readErr := updateIndexService.objectService.GetObjectSize(blobHash)
		if readErr != nil {
			return readErr
		}

		newEntry, err := domain.NewIndexEntry(p,
			blobHash,
			size,
			domain.ParseFileModeFromOsMode(fileStatInfo.Mode).Uint32(),
			fileStatInfo.Device,
			fileStatInfo.Inode,
			fileStatInfo.UserId,
			fileStatInfo.GroupId,
			domain.ComputeIndexFlags(p, 0),
			fileStatInfo.CreatedTime,
			fileStatInfo.UpdatedTime)

		if err != nil {
			return err
		}

		index.AddOrUpdateEntry(newEntry)
	}

	indexBytes, err := index.Serialize()
	if err != nil {
		return err
	}
	index.Checksum = encoding.ComputeSha256(indexBytes)

	writeErr := updateIndexService.indexService.Write(index)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func (updateIndexService *UpdateIndexService) updateIndexWithRemove(index *domain.Index, paths []string) error {
	for _, p := range paths {
		index.RemoveEntry(p)
	}
	indexBytes, err := index.Serialize()
	if err != nil {
		return err
	}
	index.Checksum = encoding.ComputeSha256(indexBytes)

	err = updateIndexService.indexService.Write(index)
	if err != nil {
		return err
	}

	return nil
}
