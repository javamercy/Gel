package workspace

import (
	"os"
)

const (
	GelDirName     = ".gel"
	ObjectsDirName = "objects"
	RefsDirName    = "refs"
	HeadFileName   = "HEAD"
	TagsDirName    = "tags"
	HeadsDirName   = "heads"
	IndexFileName  = "index"
	ConfigFileName = "config.toml"
)

const (
	FilePermission os.FileMode = 0o644 // -rw-r--r--
	DirPermission  os.FileMode = 0o755 // drwxr-xr-x
)

const (
	DefaultHeadContent = "ref: refs/heads/main\n"
)
