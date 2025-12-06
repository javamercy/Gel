package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
	"time"
)

type IUpdateIndexService interface {
	UpdateIndex(request *dto.UpdateIndexRequest) *gelErrors.GelError
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

func (updateIndexService *UpdateIndexService) UpdateIndex(request *dto.UpdateIndexRequest) *gelErrors.GelError {
	validator := validators.NewUpdateIndexValidator()
	gelError := validator.Validate(request)
	if gelError != nil {
		return gelError
	}

	err := utilities.RunAll(
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

	hashObjectRequest := dto.NewHashObjectRequest(paths, objects.GelBlobObjectType, true)
	hashMap, err := updateIndexService.hashObjectService.HashObject(hashObjectRequest)
	if err != nil {
		return err
	}

	for _, path := range paths {

		fileStatInfo, err := utilities.GetFileStatFromPath(path)
		if err != nil {
			return gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
		}

		newEntry := domain.NewIndexEntry(path,
			hashMap[path],
			fileStatInfo.Size,
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
	index.Checksum = encoding.ComputeHash(indexBytes)

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
	index.Checksum = encoding.ComputeHash(indexBytes)

	err := updateIndexService.indexRepository.Write(index)
	if err != nil {
		return gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}
	return nil
}
