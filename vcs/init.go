package vcs

import (
	"Gel/core/constant"
	"fmt"
	"path/filepath"
)

type InitService struct {
	filesystemService *FilesystemService
}

func NewInitService(filesystemService *FilesystemService) *InitService {
	return &InitService{
		filesystemService: filesystemService,
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
		filepath.Join(base, constant.GelHeadSymlinkName),
	}

	for _, dir := range dirs {
		if err := initService.filesystemService.MakeDirectory(
			dir,
			constant.GelDirPermission); err != nil {
			return "", err
		}
	}
	for _, file := range files {
		if err := initService.filesystemService.WriteFile(
			file,
			[]byte{}, false,
			constant.GelFilePermission); err != nil {
			return "", err
		}
	}

	exists := initService.filesystemService.Exists(base)
	if exists {
		return "Reinitialized existing Gel repository", nil
	}
	return fmt.Sprintf("Initialized empty Gel repository in %v", base), nil
}
