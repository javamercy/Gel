package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
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

// HashObjects hashes the contents of the files at the given paths.
// If write is true, it writes the hashed objects to the object storage.
// It returns a map of file paths to their corresponding hashes.
func (hashObjectService *HashObjectService) HashObjects(writer io.Writer, paths []string, write bool) error {

	for _, path := range paths {
		hash, serializedData, err := hashObjectService.HashObject(path)
		if err != nil {
			return err
		}

		if write {
			if err := hashObjectService.objectService.Write(hash, serializedData); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(writer, hash); err != nil {
			return err
		}
	}
	return nil
}

func (hashObjectService *HashObjectService) HashObject(path string) (string, []byte, error) {

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

	return hash, serializedData, nil
}
