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

func (h *HashObjectService) HashObjects(paths []string, options HashObjectOptions) (map[string]string, error) {
	hashes := make(map[string]string, len(paths))
	for _, path := range paths {
		hash, err := h.HashObject(path, options)
		if err != nil {
			return nil, err
		}
		hashes[path] = hash
	}
	return hashes, nil
}

func (h *HashObjectService) HashObject(path string, options HashObjectOptions) (string, error) {
	hash, serializedData, err := h.ComputeObjectHash(path)
	if err != nil {
		return "", err
	}
	if options.Write {
		if err := h.objectService.Write(hash, serializedData); err != nil {
			return "", fmt.Errorf("hash object: failed to write object to database: %w", err)
		}
	}
	return hash, nil
}

func (h *HashObjectService) ComputeObjectHash(path string) (string, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("hash object: failed to read file at '%s': %w", path, err)
	}

	blob := domain.NewBlob(data)
	serializedData := blob.Serialize()
	hash := ComputeSHA256(serializedData)
	return hash, serializedData, nil
}
