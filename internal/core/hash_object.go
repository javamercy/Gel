package core

import (
	"Gel/domain"
	"fmt"
	"io"
	"os"
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

func (h *HashObjectService) HashObjectsAndOutput(writer io.Writer, paths []string, options HashObjectOptions) error {
	for _, path := range paths {
		if err := h.HashObjectAndOutput(writer, path, options); err != nil {
			return err
		}
	}
	return nil
}

func (h *HashObjectService) HashObjectAndOutput(writer io.Writer, path string, options HashObjectOptions) error {
	hash, err := h.HashObject(path, options)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(writer, "%s\n", hash); err != nil {
		return fmt.Errorf("hash object: failed to write hash to writer: %w", err)
	}
	return nil
}

func (h *HashObjectService) HashObjects(paths []string, options HashObjectOptions) (map[string]domain.Hash, error) {
	hashes := make(map[string]domain.Hash, len(paths))
	for _, path := range paths {
		hash, err := h.HashObject(path, options)
		if err != nil {
			return nil, err
		}
		hashes[path] = hash
	}
	return hashes, nil
}

func (h *HashObjectService) HashObject(path string, options HashObjectOptions) (domain.Hash, error) {
	hash, serializedData, err := h.ComputeObjectHash(path)
	if err != nil {
		return domain.Hash{}, err
	}
	if options.Write {
		if err := h.objectService.Write(hash, serializedData); err != nil {
			return domain.Hash{}, fmt.Errorf("hash object: failed to write object to database: %w", err)
		}
	}
	return hash, nil
}

func (h *HashObjectService) ComputeObjectHash(path string) (domain.Hash, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Hash{}, nil, fmt.Errorf("hash object: failed to read file at '%s': %w", path, err)
	}

	blob := domain.NewBlob(data)
	serializedData := blob.Serialize()
	hexHash := ComputeSHA256(serializedData)
	hash, err := domain.NewHash(hexHash)
	if err != nil {
		return domain.Hash{}, nil, fmt.Errorf("hash object: failed to compute hash: %w", err)
	}
	return hash, serializedData, nil
}
