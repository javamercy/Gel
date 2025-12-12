package constant

import "os"

// repository related constants
const (
	GelRepositoryName       = ".gel"
	GelObjectsDirectoryName = "objects"
	GelRefsDirectoryName    = "refs"
	GelIndexFileName        = "index"
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
	TabStr      = "\t"
	TabByte     = '\t'
)

// permission constants
const (
	GelFilePermission      os.FileMode = 0644 // -rw-r--r----r--
	GelDirectoryPermission os.FileMode = 0755 // wxr-xr-x
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
