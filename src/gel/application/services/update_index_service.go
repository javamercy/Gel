package services

import (
	"Gel/src/gel/application/rules"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/domain"
	"Gel/src/gel/persistence/repositories"
)

type UpdateIndexOptions struct {
	Add    bool
	Remove bool
}

type IUpdateIndexService interface {
	UpdateIndex(paths []string, options UpdateIndexOptions) error
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

func (updateIndexService *UpdateIndexService) UpdateIndex(paths []string, options UpdateIndexOptions) error {

	err := updateIndexService.updateIndexRules.AllPathsMustExist(paths)
	if err != nil {
		return err
	}

	err = updateIndexService.updateIndexRules.NoDuplicatePaths(paths)
	if err != nil {
		return err
	}

	err = updateIndexService.updateIndexRules.PathsMustBeFiles(paths)
	if err != nil {
		return err
	}

	// Try to read existing index, or create new one if it doesn't exist
	index, err := updateIndexService.indexRepository.Read()
	if err != nil {
		// If index doesn't exist, create a new empty one
		index = domain.NewEmptyIndex()
	}

	if options.Add {
		err := updateIndexService.add(index, paths)
		if err != nil {
			return err
		}
	} else if options.Remove {
		err := updateIndexService.remove(index, paths)
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

		// Create blob object and get hash using HashObjectService
		// This will write the blob to the object store
		hash, err := updateIndexService.hashObjectService.HashObject(path, constant.Blob, true)
		if err != nil {
			return err
		}

		// Create or update index entry
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

		// Use domain method to add or update entry
		index.AddOrUpdateEntry(newEntry)
	}

	// Calculate checksum for the entire index
	indexBytes := serialization.SerializeIndex(index)
	index.Checksum = encoding.ComputeHash(indexBytes)

	// Write updated index to disk
	return updateIndexService.indexRepository.Write(index)

}

func (updateIndexService *UpdateIndexService) remove(index *domain.Index, paths []string) error {
	// Remove entries from index
	for _, path := range paths {
		index.RemoveEntry(path)
	}

	// Calculate checksum for the updated index
	indexBytes := serialization.SerializeIndex(index)
	index.Checksum = encoding.ComputeHash(indexBytes)

	// Write updated index to disk
	return updateIndexService.indexRepository.Write(index)
}
