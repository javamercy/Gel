package inspect

import "errors"

var (
	// ErrObjectNotFound is returned when a requested object hash does not exist in the store.
	ErrObjectNotFound = errors.New("object not found")

	// ErrUnsupportedObjectType is returned when an object type has no handler.
	ErrUnsupportedObjectType = errors.New("unsupported object type")

	// ErrInvalidRestoreMode is returned when an unknown RestoreMode is passed.
	ErrInvalidRestoreMode = errors.New("invalid restore mode")
)
