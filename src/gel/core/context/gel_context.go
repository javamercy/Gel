package context

import (
	"Gel/src/gel/core/constant"
	"os"
	"path/filepath"
	"sync"
)

type GelContext struct {
	gelDir     string
	objectsDir string
	refsDir    string
	indexPath  string
	workingDir string
}

var globalContext *GelContext
var contextOnce sync.Once

func InitializeContext() (*GelContext, error) {
	var err error
	contextOnce.Do(func() {
		cwd, e := os.Getwd()
		if e != nil {
			err = e
			return
		}

		gelDir, e := findGelDir(cwd)
		if e != nil {
			err = e
			return
		}

		globalContext = &GelContext{
			gelDir:     gelDir,
			objectsDir: filepath.Join(gelDir, constant.ObjectsDirName),
			refsDir:    filepath.Join(gelDir, constant.RefsDirName),
			indexPath:  filepath.Join(gelDir, constant.IndexFileName),
			workingDir: filepath.Dir(gelDir),
		}
	})
	return globalContext, err
}

func GetContext() *GelContext {
	if globalContext == nil {
		panic("GelContext not initialized")
	}
	return globalContext
}

func findGelDir(startPath string) (string, error) {
	currentPath := startPath
	for {
		gelPath := filepath.Join(currentPath, constant.GelDirName)
		info, err := os.Stat(gelPath)
		if err == nil && info.IsDir() {
			return gelPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return "", os.ErrNotExist
}
