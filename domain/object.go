package domain

import (
	"Gel/core/constant"
	"Gel/core/util"
	"errors"
	"strconv"
)

var (
	ErrNoNullByteFound    = errors.New("invalid object format: header must be terminated with null byte")
	ErrObjectSizeMismatch = errors.New("invalid object format: data size does not match header size")
	ErrNoSpaceInHeader    = errors.New("invalid object header: type and size must be separated by space")
	ErrUnknownObjectType  = errors.New("invalid object header: unknown object type (expected 'blob' or 'tree')")
	ErrInvalidSizeFormat  = errors.New("invalid object header: size must be a valid integer")
)

type IObject interface {
	Type() ObjectType
	Size() int
	Serialize() []byte
	Body() []byte
}

var (
	_ IObject = (*Blob)(nil)
	_ IObject = (*Tree)(nil)
	_ IObject = (*Commit)(nil)
)

func DeserializeObject(data []byte) (IObject, error) {
	nullIndex := util.FindNullByteIndex(data)
	if nullIndex == -1 {
		return nil, ErrNoNullByteFound
	}

	header := data[:nullIndex]
	objectType, size, err := deserializeObjectHeader(header)
	if err != nil {
		return nil, err
	}

	body := data[nullIndex+1:]
	if len(body) != size {
		return nil, ErrObjectSizeMismatch
	}

	switch objectType {
	case ObjectTypeBlob:
		return NewBlob(body)

	case ObjectTypeTree:
		return NewTree(body)

	case ObjectTypeCommit:
		return NewCommit(body)
	}
	// code will never reach here due to earlier validation
	return nil, nil
}

func deserializeObjectHeader(data []byte) (ObjectType, int, error) {
	spaceIndex := -1
	for i, b := range data {
		if b == constant.SpaceByte {
			spaceIndex = i
			break
		}
	}
	if spaceIndex == -1 {
		return "", 0, ErrNoSpaceInHeader
	}

	objectTypeStr := string(data[:spaceIndex])
	objectType, valid := ParseObjectType(objectTypeStr)
	if !valid {
		return "", 0, ErrUnknownObjectType
	}

	sizeStr := string(data[spaceIndex+1:])
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", 0, ErrInvalidSizeFormat
	}

	return objectType, size, nil
}

func SerializeObject(objectType ObjectType, body []byte) []byte {
	header := string(objectType) + constant.SpaceStr +
		strconv.Itoa(len(body)) + constant.NullStr
	return append([]byte(header), body...)
}
