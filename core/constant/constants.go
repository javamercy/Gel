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
	TabStr      = "\t"
	TabByte     = '\t'
)

// permission constants
const (
	GelFilePermission      os.FileMode = 0644 // -rw-r--r----r--
	GelDirectoryPermission os.FileMode = 0755 // wxr-xr-x
)

// mode constants
const (
	GelRegularFileModeStr      = "100644"
	GelRegularFileModeOctal    = 0o100644
	GelExecutableFileModeStr   = "100755"
	GelExecutableFileModeOctal = 0o100755
	GelDirectoryModeStr        = "040000"
	GelDirectoryModeOctal      = 0o040000
	GelSymlinkModeStr          = "120000"
	GelSymlinkModeOctal        = 0o120000
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
