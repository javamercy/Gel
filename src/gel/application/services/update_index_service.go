package services

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/application/rules"
	"Gel/src/gel/application/validators"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/serialization"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain"
	"Gel/src/gel/persistence/repositories"
	"errors"
	"syscall"
)

type IUpdateIndexService interface {
	UpdateIndex(request *dto.UpdateIndexRequest) error
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

func (updateIndexService *UpdateIndexService) UpdateIndex(request *dto.UpdateIndexRequest) error {
	validator := validators.NewUpdateIndexValidator()
	validationResult := validator.Validate(request)

	if !validationResult.IsValid() {
		return errors.New(validationResult.Error())
	}

	err := utilities.RunAll(
		updateIndexService.updateIndexRules.PathsMustNotDuplicate(request.Paths))

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

	hashObjectRequest := dto.NewHashObjectRequest(paths, constant.GelBlobObjectType, true)
	hashMap, err := updateIndexService.hashObjectService.HashObject(hashObjectRequest)
	if err != nil {
		return err
	}

	for _, path := range paths {
		fileInfo, err := updateIndexService.filesystemRepository.Stat(path)
		if err != nil {
			return err
		}

		statInfo, ok := fileInfo.Sys().(*syscall.Stat_t)

		if !ok {
			return errors.New("failed to get file stat info")
		}

		device, inode, userId, groupId := getFileStatSysInfo(statInfo)

		newEntry := domain.NewIndexEntry(path,
			hashMap[path],
			uint32(fileInfo.Size()),
			uint32(fileInfo.Mode()),
			device,
			inode,
			userId,
			groupId,
			getIndexFlags(path, 11),
			fileInfo.ModTime(),
			fileInfo.ModTime())

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

func getFileStatSysInfo(fileInfo *syscall.Stat_t) (uint32, uint32, uint32, uint32) {
	device := uint32(fileInfo.Dev)
	inode := uint32(fileInfo.Ino)
	userId := fileInfo.Uid
	groupId := fileInfo.Gid
	return device, inode, userId, groupId
}

func getIndexFlags(path string, stage uint16) uint16 {
	pathLength := min(len(path), 0xFFF)
	flags := uint16(pathLength) | (stage << 12)
	return flags
}
