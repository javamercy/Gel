package services

import (
	"Gel/application/dto"
	"Gel/application/rules"
	"Gel/application/validators"
	"Gel/core/crossCuttingConcerns/gelErrors"
	encoding2 "Gel/core/encoding"
	"Gel/core/utilities"
	"Gel/domain/objects"
	repositories2 "Gel/persistence/repositories"
)

type IHashObjectService interface {
	HashObject(request *dto.HashObjectRequest) (map[string]string, *gelErrors.GelError)
}

type HashObjectService struct {
	filesystemRepository repositories2.IFilesystemRepository
	objectRepository     repositories2.IObjectRepository
	hashObjectRules      *rules.HashObjectRules
}

func NewHashObjectService(filesystemRepository repositories2.IFilesystemRepository,
	objectRepository repositories2.IObjectRepository, hashObjectRules *rules.HashObjectRules) *HashObjectService {
	return &HashObjectService{
		filesystemRepository,
		objectRepository,
		hashObjectRules,
	}
}

func (hashObjectService *HashObjectService) HashObject(request *dto.HashObjectRequest) (map[string]string, *gelErrors.GelError) {

	validator := validators.NewHashObjectValidator()
	gelError := validator.Validate(request)
	if gelError != nil {
		return nil, gelError
	}

	err := utilities.Run(
		hashObjectService.hashObjectRules.PathsMustBeFiles(request.Paths),
		hashObjectService.hashObjectRules.AllPathsMustExist(request.Paths))

	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	hashMap, contentMap, err := hashObjectService.hashObjects(request.Paths)

	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}
	if !request.Write {
		return hashMap, nil
	}

	err = hashObjectService.write(contentMap)
	if err != nil {
		return nil, gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	return hashMap, nil
}

func (hashObjectService *HashObjectService) hashObjects(paths []string) (map[string]string, map[string][]byte, error) {
	hashMap := make(map[string]string)
	contentMap := make(map[string][]byte)
	for _, path := range paths {
		data, err := hashObjectService.filesystemRepository.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		blob := objects.NewBlob(data)
		content := blob.Serialize()
		hash := encoding2.ComputeHash(content)
		hashMap[path] = hash
		contentMap[hash] = content
	}
	return hashMap, contentMap, nil
}

func (hashObjectService *HashObjectService) write(contentMap map[string][]byte) error {
	for hash, content := range contentMap {
		compressedContent, err := encoding2.Compress(content)
		if err != nil {
			return err
		}

		writeErr := hashObjectService.objectRepository.Write(hash, compressedContent)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
