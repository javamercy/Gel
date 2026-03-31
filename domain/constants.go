package domain

import "os"

// Hash constants for SHA-256
const (
	SHA256HexLength  = 64 // Length of SHA-256 hash in hexadecimal
	SHA256ByteLength = 32 // Length of SHA-256 hash in bytes
)

// Index format constants
const (
	IndexSignature = "DIRC" // DirectoryMode Cache signature
	IndexVersion   = 2      // Index format version
)

const (
	GelDirName     string = ".gel"
	ObjectsDirName string = "objects"
	RefsDirName    string = "refs"
	HeadFileName   string = "HEAD"
	HeadsDirName   string = "heads"
	IndexFileName  string = "index"
	ConfigFileName string = "config.toml"
	MainBranchName string = "main"
)

const (
	MainRef string = "refs/heads/main"
)

const (
	FilePermission os.FileMode = 0o644 // -rw-r--r--
	DirPermission  os.FileMode = 0o755 // drwxr-xr-x
)
