package services

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain"
	"Gel/src/gel/persistence/repositories"
)

type UpdateIndexRequest struct {
	Paths  []string
	Add    bool
	Remove bool
}

type IUpdateIndexService interface {
	UpdateIndex(request UpdateIndexRequest) error
}

type UpdateIndexService struct {
	indexRepository      repositories.IIndexRepository
	filesystemRepository repositories.IFilesystemRepository
	hashObjectService    IHashObjectService
	updateIndexRules     *rules.UpdateIndexRules
}

func NewUpdateIndexService(indexRepository repositories.IIndexRepository, filesystemRepository repositories.IFilesystemRepository, hashObjectService IHashObjectService, updateIndexRules *rules.UpdateIndexRules) *UpdateIndexService {
	return &UpdateIndexService{
		indexRepository,
		filesystemRepository,
		hashObjectService,
		updateIndexRules,
	}
}

func (updateIndexService *UpdateIndexService) UpdateIndex(request UpdateIndexRequest) error {

	err := utilities.RunAll(
		updateIndexService.updateIndexRules.AllPathsMustExist(request.Paths),
		updateIndexService.updateIndexRules.NoDuplicatePaths(request.Paths),
		updateIndexService.updateIndexRules.PathsMustBeFiles(request.Paths))

	if err != nil {
		return err
	}

	index, err := updateIndexService.indexRepository.Read()
	if err != nil {
		index = domain.NewEmptyIndex()
	}

	if request.Add {
		err := updateIndexService.add(index, request.Paths)
		if err != nil {
			return err
		}
	} else if request.Remove {
		err := updateIndexService.remove(index, request.Paths)
		if err != nil {
			return err
		}
	}
	return nil
}

func (updateIndexService *UpdateIndexService) add(index *domain.Index, paths []string) error {

	for _, path := range paths {
		fileInfo, err := updateIndexService.filesystemRepository.Stat(path)
		if err != nil {
			return err
		}

		hash, err := updateIndexService.hashObjectService.HashObject(path, constant.Blob, true)
		if err != nil {
			return err
		}

		newEntry := domain.IndexEntry{
			Path:        path,
			Hash:        hash,
			Size:        uint32(fileInfo.Size()),
			Mode:        uint32(fileInfo.Mode()),
			Device:      0,
			Inode:       0,
			UserId:      0,
			GroupId:     0,
			Flags:       0,
			CreatedTime: fileInfo.ModTime(),
			UpdatedTime: fileInfo.ModTime(),
		}

		index.AddOrUpdateEntry(newEntry)
	}

	indexBytes := serialization.SerializeIndex(index)
	index.Checksum = encoding.ComputeHash(indexBytes)

	return updateIndexService.indexRepository.Write(index)
}

func (updateIndexService *UpdateIndexService) remove(index *domain.Index, paths []string) error {
	for _, path := range paths {
		index.RemoveEntry(path)
	}

	indexBytes := serialization.SerializeIndex(index)
	index.Checksum = encoding.ComputeHash(indexBytes)

	return updateIndexService.indexRepository.Write(index)
}
