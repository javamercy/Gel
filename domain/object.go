package domain

import (
	"strconv"
)

type Object interface {
	Type() ObjectType
	Size() int
	Serialize() []byte
	Body() []byte
}

func DeserializeObject(data []byte) (Object, error) {
	nullIndex := FindNullByteIndex(data)
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
		return NewBlob(body), nil

	case ObjectTypeTree:
		return NewTree(body)

	case ObjectTypeCommit:
		return NewCommit(body)

	default:
		return nil, ErrUnknownObjectType
	}
}

func deserializeObjectHeader(data []byte) (ObjectType, int, error) {
	spaceIndex := -1
	for i, b := range data {
		if b == ' ' {
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
	header := string(objectType) + " " +
		strconv.Itoa(len(body)) + "\x00"
	return append([]byte(header), body...)
}
