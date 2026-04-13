package domain

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	// ErrNoNullByteFound is returned when an object header is not null-terminated.
	ErrNoNullByteFound = errors.New("invalid object format: header must be terminated with null byte")

	// ErrObjectSizeMismatch is returned when the body length does not match the header size.
	ErrObjectSizeMismatch = errors.New("invalid object format: data size does not match header size")

	// ErrNoSpaceInHeader is returned when the object header does not contain the required separator.
	ErrNoSpaceInHeader = errors.New("invalid object header: type and size must be separated by space")

	// ErrUnknownObjectType is returned when the object type is not supported by the domain.
	ErrUnknownObjectType = errors.New("invalid object header: unknown object type (expected 'blob' or 'tree')")

	// ErrInvalidSizeFormat is returned when the object size field is not a valid integer.
	ErrInvalidSizeFormat = errors.New("invalid object header: size must be a valid integer")
)

// Object is implemented by Blob, Tree, and Commit.
type Object interface {
	// Type returns the object type (blob, tree, or commit).
	Type() ObjectType

	// Size returns the byte length of the object's body.
	Size() int

	// Serialize returns the full serialization in "<type> <size>\x00<body>" format.
	Serialize() []byte

	// Body returns a defensive copy of the raw object body bytes.
	Body() []byte
}

// DeserializeObject parses a serialized object from raw bytes.
// It validates the header format, checks size consistency, and dispatches
// to the appropriate constructor based on object type.
func DeserializeObject(data []byte) (Object, error) {
	dataCopy := append([]byte(nil), data...)
	nullIndex := bytes.IndexByte(dataCopy, 0)
	if nullIndex == -1 {
		return nil, ErrNoNullByteFound
	}

	header := dataCopy[:nullIndex]
	objectType, size, err := deserializeObjectHeader(header)
	if err != nil {
		return nil, err
	}

	body := dataCopy[nullIndex+1:]
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

// SerializeObject returns the full object serialization in the form
// "<type> <size>\x00<body>".
func SerializeObject(objectType ObjectType, body []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString(objectType.String())
	buf.WriteByte(' ')
	buf.WriteString(strconv.Itoa(len(body)))
	buf.WriteByte(0)
	buf.Write(body)
	return buf.Bytes()
}

// deserializeObjectHeader parses the object header ("<type> <size>") from raw bytes
// and returns the object type and body size.
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
