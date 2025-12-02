package objects

import (
	"Gel/src/gel/core/constant"
	"errors"
	"strconv"
)

type IObject interface {
	Type() ObjectType
	Size() int
	Data() []byte
	Serialize() []byte
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

func DeserializeObject(content []byte) (IObject, error) {
	nullIndex := -1
	for i, b := range content {
		if b == constant.NullByte {
			nullIndex = i
			break
		}
	}
	if nullIndex == -1 {
		return nil, errors.New("invalid object format: no null byte found")
	}

	header := content[:nullIndex]
	objectType, size, err := deserializeHeader(header)
	if err != nil {
		return nil, err
	}

	data := content[nullIndex+1:]
	if len(data) != size {
		return nil, errors.New("invalid object format: size mismatch")
	}

	switch objectType {
	case GelBlobObjectType:
		return NewBlob(data), nil
	}

	// code will never reach here due to earlier validation
	return nil, nil
}

func deserializeHeader(data []byte) (ObjectType, int, error) {
	spaceIndex := -1
	for i, b := range data {
		if b == constant.SpaceByte {
			spaceIndex = i
			break
		}
	}
	if spaceIndex == -1 {
		return "", 0, errors.New("invalid header format: no space found")
	}

	objectTypeStr := string(data[:spaceIndex])
	objectType, valid := ParseObjectType(objectTypeStr)
	if !valid {
		return "", 0, errors.New("invalid header format: unknown object type")
	}

	sizeStr := string(data[spaceIndex+1:])
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", 0, errors.New("invalid header format: size is not a valid integer")
	}

	return objectType, size, nil
}
