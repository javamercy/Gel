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
	NullStr     = "\x00"
	NullByte    = '\x00'
	SpaceStr    = " "
	SpaceByte   = ' '
	NewLineStr  = "\n"
	NewLineByte = '\n'
	SlashStr    = "/"
	SlashByte   = '/'
)

// permission constants
const (
	GelFilePermission os.FileMode = 0644 // -rw-r--r--
	GelDirPermission  os.FileMode = 0755 // wxr-xr-x
)

// tree modes
const (
	GelRegularFileMode = "100644"
	GelExecFileMode    = "100755"
	GelDirMode         = "040000"
)

// index constants
const (
	GelIndexSignature = "DIRC"
	GelIndexVersion   = 1
)

const (
	SHA256HexLength  = 64
	Sha256ByteLength = 32
)
