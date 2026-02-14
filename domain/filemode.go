package domain

import "errors"

var ErrInvalidFileMode = errors.New("invalid file mode")

type FileMode uint32

const (
	RegularFileMode    FileMode = 0o100644
	ExecutableFileMode FileMode = 0o100755
	DirectoryMode      FileMode = 0o040000
	InvalidMode        FileMode = 0
)

const (
	regularFileStr    = "100644"
	executableFileStr = "100755"
	directoryStr      = "40000"
)

func ParseFileMode(mode uint32) FileMode {
	switch mode {
	case uint32(RegularFileMode):
		return RegularFileMode
	case uint32(ExecutableFileMode):
		return ExecutableFileMode
	case uint32(DirectoryMode):
		return DirectoryMode
	default:
		return InvalidMode
	}
}

func ParseFileModeFromString(modeStr string) FileMode {
	switch modeStr {
	case regularFileStr:
		return RegularFileMode
	case executableFileStr:
		return ExecutableFileMode
	case directoryStr:
		return DirectoryMode
	default:
		return InvalidMode
	}
}

func ParseFileModeFromOsMode(osMode uint32) FileMode {
	if osMode&0o170000 == 0o040000 {
		return DirectoryMode
	} else if osMode&0o111 != 0 {
		return ExecutableFileMode
	}
	return RegularFileMode
}

func (f FileMode) String() string {
	switch f {
	case RegularFileMode:
		return regularFileStr
	case ExecutableFileMode:
		return executableFileStr
	case DirectoryMode:
		return directoryStr
	default:
		return ""
	}
}

func (f FileMode) Uint32() uint32 {
	return uint32(f)
}

func (f FileMode) IsValid() bool {
	return f != InvalidMode
}

func (f FileMode) IsDirectory() bool {
	return f == DirectoryMode
}

func (f FileMode) IsRegularFile() bool {
	return f == RegularFileMode
}

func (f FileMode) IsExecutableFile() bool {
	return f == ExecutableFileMode
}

func (f FileMode) ObjectType() (ObjectType, error) {
	switch f {
	case RegularFileMode, ExecutableFileMode:
		return ObjectTypeBlob, nil
	case DirectoryMode:
		return ObjectTypeTree, nil
	default:
		return "", ErrInvalidFileMode
	}
}

func (f FileMode) Equals(other FileMode) bool {
	return f == other
}
