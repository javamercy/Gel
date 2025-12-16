package domain

import (
	"Gel/core/constant"
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
	Data() []byte
	Serialize() []byte
	IsBlob() bool
	IsTree() bool
	IsCommit() bool
}
type BaseObject struct {
	objectType ObjectType
	data       []byte
}

func (baseObject *BaseObject) Type() ObjectType {
	return baseObject.objectType
}

func (baseObject *BaseObject) Size() int {
	return len(baseObject.data)
}

func (baseObject *BaseObject) Data() []byte {
	return baseObject.data
}

func (baseObject *BaseObject) Serialize() []byte {
	header := string(baseObject.objectType) + constant.SpaceStr + strconv.Itoa(baseObject.Size()) + constant.NullStr
	return append([]byte(header), baseObject.data...)
}

func (baseObject *BaseObject) IsBlob() bool {
	return baseObject.objectType == ObjectTypeBlob
}

func (baseObject *BaseObject) IsTree() bool {
	return baseObject.objectType == ObjectTypeTree
}

func (baseObject *BaseObject) IsCommit() bool {
	return baseObject.objectType == ObjectTypeCommit
}

func DeserializeObject(content []byte) (IObject, error) {
	nullIndex := -1
	for i, b := range content {
		if b == constant.NullByte {
			nullIndex = i
			break
		}
	}
	if nullIndex == -1 {
		return nil, ErrNoNullByteFound
	}

	header := content[:nullIndex]
	objectType, size, err := deserializeObjectHeader(header)
	if err != nil {
		return nil, err
	}

	data := content[nullIndex+1:]
	if len(data) != size {
		return nil, ErrObjectSizeMismatch
	}

	switch objectType {
	case ObjectTypeBlob:
		return NewBlob(data), nil

	case ObjectTypeTree:
		return NewTree(data), nil
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
