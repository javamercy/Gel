package gel

import (
	"Gel/domain"
	"Gel/storage"
	"fmt"
	"io"
)

type HashObjectService struct {
	objectService     *ObjectService
	filesystemStorage *storage.FilesystemStorage
}

func NewHashObjectService(objectService *ObjectService, filesystemStorage *storage.FilesystemStorage) *HashObjectService {
	return &HashObjectService{
		objectService:     objectService,
		filesystemStorage: filesystemStorage,
	}
}

func (h *HashObjectService) HashObjects(writer io.Writer, paths []string, write bool) error {

	for _, path := range paths {
		hash, _, err := h.HashObject(path, write)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "%s\n", hash); err != nil {
			return err
		}
	}
	return nil
}

func (h *HashObjectService) HashObject(path string, write bool) (string, []byte, error) {

	data, err := h.filesystemStorage.ReadFile(path)
	if err != nil {
		return "", nil, err
	}

	blob, err := domain.NewBlob(data)
	if err != nil {
		return "", nil, err
	}

	serializedData := blob.Serialize()
	hash := ComputeSHA256(serializedData)

	if write {
		if err := h.objectService.Write(hash, serializedData); err != nil {
			return "", nil, err
		}
	}

	return hash, serializedData, nil
}
