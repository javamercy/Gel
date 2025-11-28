package constant

import "os"

// ObjectType represents the type of Git object
type ObjectType string

const (
	GelBlobObjectType   ObjectType = "blob"
	GelTreeObjectType   ObjectType = "tree"
	GelCommitObjectType ObjectType = "commit"
)

// repository related constants
const (
	GelDirName        = ".gel"
	GelObjectsDirName = "objects"
	GelRefsDirName    = "refs"
	GelIndexFileName  = "index"
)

// special character constants
const (
	NullByte = "\x00"
	Space    = " "
	NewLine  = "\n"
)

// permission constants
const (
	GelFilePermission os.FileMode = 0644 // -rw-r--r--
	GelDirPermission  os.FileMode = 0755 // wxr-xr-x
)

// tree modes
const (
	RegularFileMode = "100644"
	ExecFileMode    = "100755"
	DirMode         = "40000"
)

// index constants
const (
	GelIndexSignature = "DIRC"
	GelIndexVersion   = 1
)

const (
	SHA256HexLength = 64
)
