package core

import (
	"Gel/internal/domain"
	"Gel/internal/validate"
	"errors"
	"fmt"
)

// HashObjectOptions controls hash-object behavior.
type HashObjectOptions struct {
	// Write persists hashed objects into .gel/objects when true.
	// When false, hashes are computed without writing to object storage.
	Write bool
}

// HashObjectService computes blob object hashes and optionally stores objects.
type HashObjectService struct {
	objectService *ObjectService
}

// NewHashObjectService creates a hash-object service bound to object storage operations.
func NewHashObjectService(objectService *ObjectService) *HashObjectService {
	return &HashObjectService{
		objectService: objectService,
	}
}

// HashObjects computes blob hashes for the provided absolute file paths.
//
// Each path must point to a regular file. The returned map is keyed by the
// original absolute path so callers can correlate outputs to inputs.
// If options.Write is true, each computed object is also written to storage.
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

// HashObject computes the blob hash for a single file path.
//
// The hash is computed from the serialized Git-style blob object
// ("blob <size>\\x00<body>"). When options.Write is true, the serialized object
// is compressed and written under .gel/objects using the computed hash.
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
