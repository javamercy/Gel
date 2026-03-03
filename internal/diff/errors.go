package diff

import "errors"

var (
	// ErrUnsupportedDiffMode is returned when an unknown DiffMode is passed to Diff.
	ErrUnsupportedDiffMode = errors.New("unsupported diff mode")
)
