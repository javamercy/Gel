package domain

import "os"

// SHA-256 related constants used throughout the domain package.
const (
	// SHA256HexLength is the length of a SHA-256 hash encoded as hexadecimal.
	SHA256HexLength = 64

	// SHA256ByteLength is the length of a SHA-256 hash in raw bytes.
	SHA256ByteLength = 32
)

// Repository layout constants used to construct paths under .gel.
const (
	// GelDirName is the repository metadata directory name.
	GelDirName string = ".gel"

	// ObjectsDirName is the object storage directory name.
	ObjectsDirName string = "objects"

	// RefsDirName is the references directory name.
	RefsDirName string = "refs"

	// HeadFileName is the symbolic HEAD reference filename.
	HeadFileName string = "HEAD"

	// HeadsDirName is the refs/heads directory name.
	HeadsDirName string = "heads"

	// IndexFileName is the index file name.
	IndexFileName string = "index"

	// ConfigFileName is the repository config file name.
	ConfigFileName string = "config.toml"

	// MainBranchName is the default branch name.
	MainBranchName string = "main"
)

// MainRef is the full ref path for the default branch.
const (
	MainRef string = "refs/heads/main"
)

// File permission constants used when writing repository files and directories.
const (
	// FilePermission is the default permission for regular files written by Gel.
	FilePermission os.FileMode = 0o644

	// DirPermission is the default permission for directories created by Gel.
	DirPermission os.FileMode = 0o755
)
