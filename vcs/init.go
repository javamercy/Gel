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
		base,
		filepath.Join(base, constant.GelObjectsDirectoryName),
		filepath.Join(base, constant.GelRefsDirectoryName),
	}

	files := []string{
		filepath.Join(base, constant.GelConfigFileName),
	}

	exists := initService.filesystemService.Exists(base)
	for _, dir := range dirs {
		if err := initService.filesystemService.MakeDirectory(
			dir,
			constant.GelDirectoryPermission); err != nil {
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

	if exists {
		return "Reinitialized existing Gel repository", nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %v", base), nil
}
