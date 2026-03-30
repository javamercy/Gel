package core

import (
	"Gel/domain"
	"Gel/internal/validate"
	"errors"
	"fmt"
)

type HashObjectOptions struct {
	Write bool
}
type HashObjectService struct {
	objectService *ObjectService
}

func NewHashObjectService(objectService *ObjectService) *HashObjectService {
	return &HashObjectService{
		objectService: objectService,
	}
}

func (h *HashObjectService) HashObjects(paths []domain.AbsolutePath, options HashObjectOptions) (
	map[domain.AbsolutePath]domain.Hash, error,
) {
	if len(paths) == 0 {
		return nil, errors.New("no paths provided")
	}

	hashes := make(map[domain.AbsolutePath]domain.Hash, len(paths))
	for _, path := range paths {
		if err := validate.PathMustBeFile(path.String()); err != nil {
			return nil, fmt.Errorf("hash-object: %w", err)
		}

		hash, err := h.HashObject(path, options)
		if err != nil {
			return nil, err
		}

		hashes[path] = hash
	}
	return hashes, nil
}

func (h *HashObjectService) HashObject(path domain.AbsolutePath, options HashObjectOptions) (domain.Hash, error) {
	hash, serializedData, err := h.objectService.ComputeObjectHash(path)
	if err != nil {
		return domain.Hash{}, err
	}
	if options.Write {
		if err := h.objectService.Write(hash, serializedData); err != nil {
			return domain.Hash{}, fmt.Errorf("hash object: %w", err)
		}
	}
	return hash, nil
}
