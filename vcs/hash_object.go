package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
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

// HashObject hashes the contents of the files at the given paths.
// If write is true, it writes the hashed objects to the object storage.
// It returns a map of file paths to their corresponding hashes.
func (hashObjectService *HashObjectService) HashObject(paths []string, write bool) (map[string]string, error) {

	hashMap := make(map[string]string)

	for _, path := range paths {
		data, err := hashObjectService.filesystemService.ReadFile(path)
		if err != nil {
			return nil, err
		}

		blob, err := domain.NewBlob(data)
		if err != nil {
			return nil, err
		}

		content := blob.Serialize()
		hash := encoding.ComputeSha256(content)
		hashMap[path] = hash

		if write {
			if err := hashObjectService.objectService.Write(hash, content); err != nil {
				return nil, err
			}
		}
	}
	return hashMap, nil
}
