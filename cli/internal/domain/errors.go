package domain

import "errors"

// ErrInvalidObjectType is returned when an object type is invalid or unsupported.
var ErrInvalidObjectType = errors.New("invalid object type")
