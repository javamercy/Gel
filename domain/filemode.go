package domain

import "errors"

type FileMode uint32

var ErrInvalidFileMode = errors.New("invalid file mode")

const (
	RegularFile    FileMode = 0o100644
	ExecutableFile FileMode = 0o100755
	Directory      FileMode = 0o040000
	Symlink        FileMode = 0o120000
	Submodule      FileMode = 0o160000
	InvalidMode    FileMode = 0
)

const (
	RegularFileStr    = "100644"
	ExecutableFileStr = "100755"
	DirectoryStr      = "40000"
	SymlinkStr        = "120000"
	SubmoduleStr      = "160000"
)

func ParseFileMode(mode uint32) FileMode {
	switch mode {
	case uint32(RegularFile):
		return RegularFile
	case uint32(ExecutableFile):
		return ExecutableFile
	case uint32(Directory):
		return Directory
	case uint32(Symlink):
		return Symlink
	case uint32(Submodule):
		return Submodule
	default:
		return InvalidMode
	}
}

func ParseFileModeFromString(modeStr string) FileMode {
	switch modeStr {
	case RegularFileStr:
		return RegularFile
	case ExecutableFileStr:
		return ExecutableFile
	case DirectoryStr:
		return Directory
	case SymlinkStr:
		return Symlink
	case SubmoduleStr:
		return Submodule
	default:
		return InvalidMode
	}
}

func ParseFileModeFromOsMode(osMode uint32) FileMode {
	if osMode&0o170000 == 0o040000 {
		return Directory
	} else if osMode&0o170000 == 0o120000 {
		return Symlink
	} else if osMode&0o170000 == 0o160000 {
		return Submodule
	} else if osMode&0o111 != 0 {
		return ExecutableFile
	}

	return RegularFile
}

func (filemode FileMode) String() string {
	switch filemode {
	case RegularFile:
		return RegularFileStr
	case ExecutableFile:
		return ExecutableFileStr
	case Directory:
		return DirectoryStr
	case Symlink:
		return SymlinkStr
	case Submodule:
		return SubmoduleStr
	default:
		return ""
	}
}

func (filemode FileMode) Uint32() uint32 {
	return uint32(filemode)
}

func (filemode FileMode) IsValid() bool {
	return filemode != InvalidMode
}

func (filemode FileMode) IsDirectory() bool {
	return filemode == Directory
}

func (filemode FileMode) IsRegularFile() bool {
	return filemode == RegularFile
}

func (filemode FileMode) IsExecutableFile() bool {
	return filemode == ExecutableFile
}

func (filemode FileMode) IsSymlink() bool {
	return filemode == Symlink
}

func (filemode FileMode) IsSubmodule() bool {
	return filemode == Submodule
}

func (filemode FileMode) ObjectType() (ObjectType, error) {
	switch filemode {
	case RegularFile, ExecutableFile, Symlink:
		return ObjectTypeBlob, nil
	case Directory:
		return ObjectTypeTree, nil
	case Submodule:
		return ObjectTypeCommit, nil
	default:
		return "", ErrInvalidFileMode
	}
}

func (filemode FileMode) Equals(other FileMode) bool {
	return filemode == other
}
