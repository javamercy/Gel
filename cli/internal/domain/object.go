package domain

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var (
	// ErrObjectHeaderMissingTerminator is returned when an object header is not null-terminated.
	ErrObjectHeaderMissingTerminator = errors.New("invalid object format: missing header terminator")

	// ErrObjectSizeMismatch is returned when the body length does not match the header size.
	ErrObjectSizeMismatch = errors.New("invalid object format: body size does not match header size")

	// ErrObjectHeaderMissingSeparator is returned when the object header does not contain the required separator.
	ErrObjectHeaderMissingSeparator = errors.New("invalid object header: missing type/size separator")

	// ErrObjectTypeUnknown is returned when the object type is not supported.
	ErrObjectTypeUnknown = errors.New("invalid object header: unknown object type")

	// ErrObjectTypeMismatch is returned when an object has a different type than expected.
	ErrObjectTypeMismatch = errors.New("object type mismatch")

	// ErrObjectSizeInvalid is returned when the object size field is not a non-negative integer.
	ErrObjectSizeInvalid = errors.New("invalid object header: size must be a non-negative integer")
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
	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, ErrObjectHeaderMissingTerminator
	}

	objectType, size, err := deserializeObjectHeader(data[:nullIndex])
	if err != nil {
		return nil, err
	}

	body := data[nullIndex+1:]
	if len(body) != size {
		return nil, fmt.Errorf("%w: header=%d actual=%d", ErrObjectSizeMismatch, size, len(body))
	}

	switch objectType {
	case ObjectTypeBlob:
		return NewBlob(body), nil
	case ObjectTypeTree:
		return NewTree(body)
	case ObjectTypeCommit:
		return NewCommit(body)
	default:
		return nil, fmt.Errorf("%w: %q", ErrObjectTypeUnknown, objectType)
	}
}

// SerializeObject returns the full object serialization in the form
// "<type> <size>\x00<body>".
//
// It does not validate objectType. Callers must pass a valid ObjectType;
// validation is performed when deserializing object data.
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
	spaceIndex := bytes.IndexByte(data, ' ')
	if spaceIndex == -1 {
		return "", 0, ErrObjectHeaderMissingSeparator
	}

	objectTypeName := string(data[:spaceIndex])
	objectType, ok := ParseObjectType(objectTypeName)
	if !ok {
		return "", 0, fmt.Errorf("%w: %q", ErrObjectTypeUnknown, objectTypeName)
	}

	sizeText := string(data[spaceIndex+1:])
	size, err := strconv.Atoi(sizeText)
	if err != nil || size < 0 {
		return "", 0, fmt.Errorf("%w: %q", ErrObjectSizeInvalid, sizeText)
	}
	return objectType, size, nil
}
