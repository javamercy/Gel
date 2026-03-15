package workspace

import (
	"os"
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
