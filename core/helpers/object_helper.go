package helpers

import (
	"Gel/core/constants"
	"Gel/domain/objects"
	"errors"
	"strconv"
	"strings"
)

// ToObjectContent prepares the object content by adding a header
func ToObjectContent(objectType constants.ObjectType, fileData []byte) []byte {
	header := string(objectType) + constants.Space + strconv.Itoa(len(fileData)) + constants.NullByte
	return append([]byte(header), fileData...)
}

// ToObject converts raw object data into an IObject
func ToObject(data []byte) (objects.IObject, error) {
	nullIndex := -1
	for i, b := range data {
		if b == 0 {
			nullIndex = i
			break
		}
	}
	if nullIndex == -1 {
		return nil, errors.New("invalid object format: no null byte found")
	}

	header := string(data[:nullIndex])
	parts := strings.Split(header, constants.Space)
	if len(parts) != 2 {
		return nil, errors.New("invalid header format")
	}

	sizeStr := parts[1]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, errors.New("invalid size in header")
	}

	// Extract content after null byte
	content := data[nullIndex+1:]
	if len(content) != size {
		return nil, errors.New("content size mismatch")
	}

	objStr := parts[0]
	objType := constants.ObjectType(objStr)
	switch objType {
	case constants.Blob:
		return objects.NewBlob(content), nil
	default:
		return nil, errors.New("unsupported object type: " + objStr)
	}
}
