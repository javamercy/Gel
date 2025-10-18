package constants

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
