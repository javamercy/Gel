package constant

import "os"

// repository related constants
const (
	GelDirName        = ".gel"
	GelObjectsDirName = "objects"
	GelRefsDirName    = "refs"
	GelIndexFileName  = "index"
)

// special character constants
const (
	NullByte    = '\x00'
	Space       = " "
	SpaceByte   = ' '
	NewLine     = "\n"
	NewLineByte = '\n'
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
