package context

import (
	"Gel/core/constant"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type GelContext struct {
	GelDir        string
	ObjectsDir    string
	RefsDir       string
	IndexPath     string
	RepositoryDir string
}

var globalContext *GelContext
var contextOnce sync.Once

func EnsureContext() error {
	_, err := initializeContext()
	if err != nil {
		return errors.New("not a gel repository (or any of the parent directories): .gel")
	}
	return nil
}

func GetContext() *GelContext {
	return globalContext
}
func initializeContext() (*GelContext, error) {
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
			GelDir:        gelDir,
			ObjectsDir:    filepath.Join(gelDir, constant.GelObjectsDirName),
			RefsDir:       filepath.Join(gelDir, constant.GelRefsDirName),
			IndexPath:     filepath.Join(gelDir, constant.GelIndexFileName),
			RepositoryDir: filepath.Dir(gelDir),
		}
	})
	return globalContext, err
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
