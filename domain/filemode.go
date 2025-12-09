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
	DirectoryStr      = "040000"
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

func (fm FileMode) String() string {
	switch fm {
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

func (fm FileMode) ToUint32() uint32 {
	return uint32(fm)
}

func (fm FileMode) IsValid() bool {
	return fm != InvalidMode
}

func (fm FileMode) IsDirectory() bool {
	return fm == Directory
}

func (fm FileMode) IsRegularFile() bool {
	return fm == RegularFile
}

func (fm FileMode) IsExecutableFile() bool {
	return fm == ExecutableFile
}

func (fm FileMode) IsSymlink() bool {
	return fm == Symlink
}

func (fm FileMode) IsSubmodule() bool {
	return fm == Submodule
}

func (fm FileMode) ObjectType() (ObjectType, error) {
	switch fm {
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

func (fm FileMode) Equals(other FileMode) bool {
	return fm == other
}
