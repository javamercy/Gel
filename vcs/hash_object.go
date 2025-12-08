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
func (hashObjectService *HashObjectService) HashObject(paths []string, write bool) (map[string]string, map[string][]byte, error) {

	hashMap := make(map[string]string)
	contentMap := make(map[string][]byte)
	for _, path := range paths {
		data, err := hashObjectService.filesystemService.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		blob := domain.NewBlob(data)
		content := blob.Serialize()
		hash := encoding.ComputeHash(content)
		hashMap[path] = hash
		contentMap[hash] = content

		if write {
			if err := hashObjectService.objectService.Write(hash, content); err != nil {
				return nil, nil, err
			}
		}
	}

	return hashMap, contentMap, nil
}
