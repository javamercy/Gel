package staging

import "errors"

var (
	// ErrPathDidNotMatch is returned when a pathspec matches no files in the index or working tree.
	ErrPathDidNotMatch = errors.New("pathspec did not match any files")
)
