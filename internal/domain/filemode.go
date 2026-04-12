package domain

import "errors"

// ErrInvalidFileMode indicates that the provided file mode is invalid or unrecognized.
var ErrInvalidFileMode = errors.New("invalid file mode")

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
	// FileModeInvalid represents an unknown or invalid mode.
	FileModeInvalid FileMode = 0
)

// ParseFileMode converts a uint32 mode value to a FileMode.
func ParseFileMode(mode uint32) FileMode {
	switch mode {
	case uint32(FileModeRegular):
		return FileModeRegular
	case uint32(FileModeExecutable):
		return FileModeExecutable
	case uint32(FileModeDirectory):
		return FileModeDirectory
	default:
		return FileModeInvalid
	}
}

// ParseFileModeFromString converts a string mode (e.g., "100644") to a FileMode.
func ParseFileModeFromString(mode string) FileMode {
	switch mode {
	case FileModeRegular.String():
		return FileModeRegular
	case FileModeExecutable.String():
		return FileModeExecutable
	case FileModeDirectory.String():
		return FileModeDirectory
	default:
		return FileModeInvalid
	}
}

// ParseFileModeFromOsMode converts an OS-level mode (from syscall) to a FileMode.
func ParseFileModeFromOsMode(osMode uint32) FileMode {
	if osMode&0o170000 == 0o040000 {
		return FileModeDirectory
	} else if osMode&0o111 != 0 {
		return FileModeExecutable
	}
	return FileModeRegular
}

// String returns the string representation of the mode (e.g., "100644").
// Returns empty string for invalid modes.
func (f FileMode) String() string {
	switch f {
	case FileModeRegular:
		return "100644"
	case FileModeExecutable:
		return "100755"
	case FileModeDirectory:
		return "40000"
	default:
		return ""
	}
}

// Uint32 returns the raw uint32 value of the mode.
func (f FileMode) Uint32() uint32 {
	return uint32(f)
}

// IsValid reports whether the mode is recognized.
func (f FileMode) IsValid() bool {
	return f != FileModeInvalid
}

// IsDirectory reports whether the mode represents a directory.
func (f FileMode) IsDirectory() bool {
	return f == FileModeDirectory
}

// IsRegularFile reports whether the mode represents a regular non-executable file.
func (f FileMode) IsRegularFile() bool {
	return f == FileModeRegular
}

// IsExecutableFile reports whether the mode represents an executable file.
func (f FileMode) IsExecutableFile() bool {
	return f == FileModeExecutable
}

// ObjectType returns the domain object type corresponding to this mode.
func (f FileMode) ObjectType() (ObjectType, error) {
	switch f {
	case FileModeRegular, FileModeExecutable:
		return ObjectTypeBlob, nil
	case FileModeDirectory:
		return ObjectTypeTree, nil
	default:
		return "", ErrInvalidFileMode
	}
}

// Equals reports whether two modes are identical.
func (f FileMode) Equals(other FileMode) bool {
	return f == other
}
