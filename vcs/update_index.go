package vcs

import (
	"Gel/core/encoding"
	"Gel/core/utilities"
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

	hashMap, _, err := updateIndexService.hashObjectService.HashObject(paths, true)
	if err != nil {
		return err
	}

	for _, p := range paths {

		fileStatInfo := utilities.GetFileStatFromPath(p)

		blobHash := hashMap[p]
		size, readErr := updateIndexService.objectService.GetObjectSize(blobHash)
		if readErr != nil {
			return readErr
		}

		newEntry := domain.NewIndexEntry(p,
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

		index.AddOrUpdateEntry(newEntry)
	}

	indexBytes := index.Serialize()
	index.Checksum = encoding.ComputeHash(indexBytes)

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
	indexBytes := index.Serialize()
	index.Checksum = encoding.ComputeHash(indexBytes)

	err := updateIndexService.indexService.Write(index)
	if err != nil {
		return err
	}

	return nil
}
