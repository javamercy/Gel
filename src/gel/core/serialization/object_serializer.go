package serialization

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/domain/objects"
	"errors"
	"strconv"
	"strings"
)

func SerializeObject(objectType constant.ObjectType, fileData []byte) []byte {
	header := string(objectType) + constant.Space + strconv.Itoa(len(fileData)) + constant.NullByte
	return append([]byte(header), fileData...)
}

func DeserializeObject(data []byte) (objects.IObject, error) {
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
	parts := strings.Split(header, constant.Space)
	if len(parts) != 2 {
		return nil, errors.New("invalid header format")
	}

	sizeStr := parts[1]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return nil, errors.New("invalid size in header")
	}

	content := data[nullIndex+1:]
	if len(content) != size {
		return nil, errors.New("content size mismatch")
	}

	objStr := parts[0]
	objType := constant.ObjectType(objStr)
	switch objType {
	case constant.GelBlobObjectType:
		return objects.NewBlob(content), nil
	default:
		return nil, errors.New("unsupported object type: " + objStr)
	}
}
