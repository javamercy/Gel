package repositories

import "os"

type IRepository interface {
	MakeDir(path string, permission os.FileMode) error
	MakeDirRange(paths []string, permission os.FileMode) error
	WriteFile(path string, data []byte, autoCreateDir bool, permission os.FileMode) error
	Exists(path string) bool
	ReadFile(path string) ([]byte, error)
	FindGelDir(startPath string) (string, error)
	FindObjectPath(hash string, startPath string) (string, error)
	FindObjectsDir(startPath string) (string, error)
}
