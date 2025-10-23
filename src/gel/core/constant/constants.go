package constant

import "os"

// ObjectType represents the type of Git object
type ObjectType string

const (
	Blob   ObjectType = "blob"
	Tree   ObjectType = "tree"
	Commit ObjectType = "commit"
)

// repository related constants
const (
	GelDirName     = ".gel"
	ObjectsDirName = "objects"
	RefsDirName    = "refs"
	IndexFileName  = "index"
)

// special character constants
const (
	NullByte = "\x00"
	Space    = " "
)

// permission constants
const (
	FilePermission os.FileMode = 0644 // -rw-r--r--
	DirPermission  os.FileMode = 0755 // wxr-xr-x
)

// tree modes
const (
	RegularFileMode = "100644"
	ExecFileMode    = "100755"
	DirMode         = "40000"
)

// index constants
const (
	IndexSignature = "GELI"
	IndexVersion   = 2
)
