package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
	"fmt"
	"io"
)

type HashObjectService struct {
	objectService     *ObjectService
	filesystemService *FilesystemService
}

func NewHashObjectService(objectService *ObjectService, filesystemService *FilesystemService) *HashObjectService {
	return &HashObjectService{
		objectService:     objectService,
		filesystemService: filesystemService,
	}
}

func (hashObjectService *HashObjectService) HashObjects(writer io.Writer, paths []string, write bool) error {

	for _, path := range paths {
		hash, _, err := hashObjectService.HashObject(path, write)
		if err != nil {
			return err
		}
		if _, err := io.WriteString(writer, fmt.Sprintf("%v\n", hash)); err != nil {
			return err
		}
	}
	return nil
}

func (hashObjectService *HashObjectService) HashObject(path string, write bool) (string, []byte, error) {

	data, err := hashObjectService.filesystemService.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	blob, err := domain.NewBlob(data)
	if err != nil {
		return "", nil, err
	}

	serializedData := blob.Serialize()
	hash := encoding.ComputeSha256(serializedData)

	if write {
		if err := hashObjectService.objectService.Write(hash, serializedData); err != nil {
			return "", nil, err
		}
	}

	return hash, serializedData, nil
}
