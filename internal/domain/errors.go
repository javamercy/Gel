package domain

import "errors"

var (
	ErrInvalidFileMode   = errors.New("invalid file mode")
	ErrInvalidObjectType = errors.New("invalid object type")

	ErrInvalidCommitFormat = errors.New("invalid commit format")
	ErrInvalidCommitField  = errors.New("invalid commit field")

	ErrNoNullByteFound    = errors.New("invalid object format: header must be terminated with null byte")
	ErrObjectSizeMismatch = errors.New("invalid object format: data size does not match header size")
	ErrNoSpaceInHeader    = errors.New("invalid object header: type and size must be separated by space")
	ErrUnknownObjectType  = errors.New("invalid object header: unknown object type (expected 'blob' or 'tree')")
	ErrInvalidSizeFormat  = errors.New("invalid object header: size must be a valid integer")

	ErrIndexTooShort         = errors.New("index file is too short: minimum 12 bytes required for header")
	ErrInvalidIndexSignature = errors.New("invalid index signature: expected 'DIRC', file may be corrupted")
	ErrTruncatedEntryData    = errors.New("index file truncated: not enough data to read all entries")
	ErrIncorrectChecksumSize = errors.New("invalid index checksum: expected 32 bytes at end of file")
	ErrChecksumMismatch      = errors.New("index checksum verification failed: file may be corrupted")
	ErrEntryDataTooShort     = errors.New("index entry is incomplete: minimum 74 bytes required")
	ErrPathNotNullTerminated = errors.New("index entry path is malformed: missing null terminator")

	ErrTreeMissingNullByte = errors.New("invalid tree format: missing null byte after name")
	ErrTreeTruncatedHash   = errors.New("invalid tree format: truncated hash")
)
