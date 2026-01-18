package constant

import "os"

// repository related constants
const (
	GelRepositoryName  = ".gel"
	GelObjectsDirName  = "objects"
	GelRefsDirName     = "refs"
	GelHeadSymlinkName = "HEAD"
	GelTagsDirName     = "tags"
	GelHeadsDirName    = "heads"
	GelIndexFileName   = "index"
	GelConfigFileName  = "config.toml"
)

// permission constants
const (
	GelFilePermission os.FileMode = 0o0644 // -rw-r--r----r--
	GelDirPermission  os.FileMode = 0o0755 // wxr-xr-x
)

// index constants
const (
	GelIndexSignature = "DIRC"
	GelIndexVersion   = 1
)

// hash constants
const (
	SHA256HexLength  = 64
	Sha256ByteLength = 32
)
