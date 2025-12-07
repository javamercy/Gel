package services

import (
	"Gel/application/dto"
	"Gel/application/rules"
	"Gel/application/validators"
	"Gel/core/crossCuttingConcerns/gelErrors"
	encoding2 "Gel/core/encoding"
	utilities2 "Gel/core/utilities"
	"Gel/domain"
	objects2 "Gel/domain/objects"
	repositories2 "Gel/persistence/repositories"
	"time"
)

type IUpdateIndexService interface {
	UpdateIndex(request *dto.UpdateIndexRequest) *gelErrors.GelError
}

type UpdateIndexService struct {
	indexRepository      repositories2.IIndexRepository
	filesystemRepository repositories2.IFilesystemRepository
	objectRepository     repositories2.IObjectRepository
	hashObjectService    IHashObjectService
	updateIndexRules     *rules.UpdateIndexRules
}

func NewUpdateIndexService(indexRepository repositories2.IIndexRepository, filesystemRepository repositories2.IFilesystemRepository, objectRepository repositories2.IObjectRepository, hashObjectService IHashObjectService, updateIndexRules *rules.UpdateIndexRules) *UpdateIndexService {
	return &UpdateIndexService{
		indexRepository,
		filesystemRepository,
		objectRepository,
		hashObjectService,
		updateIndexRules,
	}
}

func (updateIndexService *UpdateIndexService) UpdateIndex(request *dto.UpdateIndexRequest) *gelErrors.GelError {
	validator := validators.NewUpdateIndexValidator()
	gelError := validator.Validate(request)
	if gelError != nil {
		return gelError
	}

	err := utilities2.RunAll(
		updateIndexService.updateIndexRules.PathsMustNotDuplicate(request.Paths))

	if err != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	index, err := updateIndexService.indexRepository.Read()
	if err != nil {
		index = domain.NewEmptyIndex()
	}

	if request.Add {
		return updateIndexService.add(index, request.Paths)

	} else if request.Remove {
		return updateIndexService.remove(index, request.Paths)

	}

	return nil
}

func (updateIndexService *UpdateIndexService) add(index *domain.Index, paths []string) *gelErrors.GelError {

	hashObjectRequest := dto.NewHashObjectRequest(paths, objects2.GelBlobObjectType, true)
	hashMap, err := updateIndexService.hashObjectService.HashObject(hashObjectRequest)
	if err != nil {
		return err
	}

	for _, path := range paths {

		fileStatInfo, err := utilities2.GetFileStatFromPath(path)
		if err != nil {
			return gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
		}

		blobHash := hashMap[path]
		size, readErr := updateIndexService.readBlobAndGetSize(blobHash)
		if readErr != nil {
			return gelErrors.NewGelError(gelErrors.ExitCodeFatal, readErr.Error())
		}

		newEntry := domain.NewIndexEntry(path,
			blobHash,
			size,
			fileStatInfo.Mode,
			fileStatInfo.Device,
			fileStatInfo.Inode,
			fileStatInfo.UserId,
			fileStatInfo.GroupId,
			domain.ComputeIndexFlags(path, 0),
			time.Now(),
			time.Now())

		index.AddOrUpdateEntry(newEntry)
	}

	indexBytes := index.Serialize()
	index.Checksum = encoding2.ComputeHash(indexBytes)

	writeErr := updateIndexService.indexRepository.Write(index)
	if writeErr != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, writeErr.Error())
	}
	return nil
}

func (updateIndexService *UpdateIndexService) remove(index *domain.Index, paths []string) *gelErrors.GelError {
	for _, path := range paths {
		index.RemoveEntry(path)
	}
	indexBytes := index.Serialize()
	index.Checksum = encoding2.ComputeHash(indexBytes)

	err := updateIndexService.indexRepository.Write(index)
	if err != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}
	return nil
}

func (updateIndexService *UpdateIndexService) readBlobAndGetSize(hash string) (uint32, error) {
	compressedContent, err := updateIndexService.objectRepository.Read(hash)
	if err != nil {
		return 0, err
	}

	content, err := encoding2.Decompress(compressedContent)
	if err != nil {
		return 0, err
	}

	object, err := objects2.DeserializeObject(content)
	if err != nil {
		return 0, err
	}

	return uint32(object.Size()), nil
}
