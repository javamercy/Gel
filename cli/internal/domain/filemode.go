package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidFileMode is returned when a mode is not one of Gel's supported file modes.
	ErrInvalidFileMode = errors.New("invalid file mode")

	// ErrUnsupportedFileType is returned when an OS file type cannot be represented by Gel.
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

// FileMode represents the type and permissions of a file in a tree entry.
type FileMode uint32

// File mode constants following Git conventions.
const (
	// FileModeRegular is a regular non-executable file (100644).
	FileModeRegular FileMode = 0o100644

	// FileModeExecutable is an executable file (100755).
	FileModeExecutable FileMode = 0o100755

	// FileModeDirectory is a directory (040000).
	FileModeDirectory FileMode = 0o040000
)

const (
	treeModeRegular             = "100644"
	treeModeExecutable          = "100755"
	treeModeDirectory           = "40000"
	expectedStoredModes         = "100644, 100755, or 040000"
	expectedTreeModes           = `"100644", "100755", or "40000"`
	osModeTypeMask       uint32 = 0o170000
	osModeTypeRegular    uint32 = 0o100000
	osModeTypeDirectory  uint32 = 0o040000
	osModeExecutableMask uint32 = 0o111
)

// NewFileMode converts a raw stored mode value to a FileMode.
func NewFileMode(mode uint32) (FileMode, error) {
	fileMode := FileMode(mode)
	if fileMode.IsValid() {
		return fileMode, nil
	}
	return 0, fmt.Errorf(
		"%w: %s (expected %s)",
		ErrInvalidFileMode,
		formatStoredFileMode(mode),
		expectedStoredModes,
	)
}

// NewFileModeFromTreeMode converts a canonical tree-mode string to a FileMode.
func NewFileModeFromTreeMode(mode string) (FileMode, error) {
	switch mode {
	case treeModeRegular:
		return FileModeRegular, nil
	case treeModeExecutable:
		return FileModeExecutable, nil
	case treeModeDirectory:
		return FileModeDirectory, nil
	default:
		return 0, fmt.Errorf(
			"%w: %q (expected %s)",
			ErrInvalidFileMode,
			mode,
			expectedTreeModes,
		)
	}
}

// NewFileModeFromOSMode converts an OS file mode to a canonical Gel file mode.
func NewFileModeFromOSMode(mode uint32) (FileMode, error) {
	switch mode & osModeTypeMask {
	case osModeTypeDirectory:
		return FileModeDirectory, nil
	case osModeTypeRegular:
		if mode&osModeExecutableMask != 0 {
			return FileModeExecutable, nil
		}
		return FileModeRegular, nil
	default:
		return 0, fmt.Errorf(
			"%w: mode=%s type=%s",
			ErrUnsupportedFileType,
			formatStoredFileMode(mode),
			formatStoredFileMode(mode&osModeTypeMask),
		)
	}
}

// String returns the canonical tree-mode string for mode.
// It returns an empty string for invalid modes.
func (f FileMode) String() string {
	switch f {
	case FileModeRegular:
		return treeModeRegular
	case FileModeExecutable:
		return treeModeExecutable
	case FileModeDirectory:
		return treeModeDirectory
	default:
		return ""
	}
}

// Uint32 returns the raw uint32 value of the mode.
func (f FileMode) Uint32() uint32 {
	return uint32(f)
}

// IsValid reports whether mode is one of Gel's supported file modes.
func (f FileMode) IsValid() bool {
	switch f {
	case FileModeRegular, FileModeExecutable, FileModeDirectory:
		return true
	default:
		return false
	}
}

// IsDirectory reports whether the mode represents a directory.
func (f FileMode) IsDirectory() bool {
	return f == FileModeDirectory
}

// ObjectType returns the domain object type corresponding to this mode.
func (f FileMode) ObjectType() (ObjectType, error) {
	switch f {
	case FileModeRegular, FileModeExecutable:
		return ObjectTypeBlob, nil
	case FileModeDirectory:
		return ObjectTypeTree, nil
	default:
		return "", fmt.Errorf(
			"%w: %s (expected %s)",
			ErrInvalidFileMode,
			formatStoredFileMode(f.Uint32()),
			expectedStoredModes,
		)
	}
}

// Equals reports whether two modes are identical.
func (f FileMode) Equals(o FileMode) bool {
	return f == o
}

// formatStoredFileMode returns a zero-padded octal string representation of mode for error messages.
func formatStoredFileMode(mode uint32) string {
	return fmt.Sprintf("%06o", mode)
}
