package constants

import "os"

// ObjectType represents the type of Git object
type ObjectType string

const (
	Blob   ObjectType = "blob"
	Tree   ObjectType = "tree"
	Commit ObjectType = "commit"
)

// Repository related constants
const (
	RepositoryDirName = ".gel"
	ObjectsDirName    = "objects"
	RefsDirName       = "refs"
)

// Special character constants
const (
	NullByte = "\x00"
	Space    = " "
)

// Permission constants
const (
	FilePermission os.FileMode = 0644 // -rw-r--r--
	DirPermission  os.FileMode = 0755 // wxr-xr-x
)
