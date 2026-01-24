package vcs

import (
	"Gel/core/constant"
	"Gel/storage"
	"fmt"
	"path/filepath"
)

type InitService struct {
	filesystemStorage *storage.FilesystemStorage
}

func NewInitService(filesystemStorage *storage.FilesystemStorage) *InitService {
	return &InitService{
		filesystemStorage: filesystemStorage,
	}
}

func (initService *InitService) Init(path string) (string, error) {

	base := filepath.Join(path, constant.GelRepositoryName)

	dirs := []string{
		filepath.Join(base, constant.GelObjectsDirName),
		filepath.Join(base, constant.GelRefsDirName, constant.GelHeadsDirName),
		filepath.Join(base, constant.GelRefsDirName, constant.GelTagsDirName),
	}

	files := []string{
		filepath.Join(base, constant.GelConfigFileName),
		filepath.Join(base, constant.GelHeadFileName),
	}

	exists := initService.filesystemStorage.Exists(base)

	for _, dir := range dirs {
		if err := initService.filesystemStorage.MakeDir(
			dir,
			constant.GelDirPermission); err != nil {
			return "", err
		}
	}
	for i, file := range files {
		var err error
		if i == 1 {
			headRefBytes := []byte("ref: refs/heads/main\n")
			err = initService.filesystemStorage.WriteFile(
				file,
				headRefBytes, false,
				constant.GelFilePermission)
		} else {
			err = initService.filesystemStorage.WriteFile(
				file,
				[]byte{}, false,
				constant.GelFilePermission)
		}
		if err != nil {
			return "", err
		}
	}

	if exists {
		return "Reinitialized existing Gel repository", nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", base), nil
}
