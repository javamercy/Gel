package vcs

import (
	"Gel/core/constant"
	"path/filepath"
)

type InitService struct {
	objectService     *ObjectService
	filesystemService *FilesystemService
}

func NewInitService(objectService *ObjectService, filesystemService *FilesystemService) *InitService {
	return &InitService{
		objectService:     objectService,
		filesystemService: filesystemService,
	}
}

func (initService *InitService) Init(path string) (string, error) {

	base := filepath.Join(path, constant.GelDirName)

	dirs := []string{
		base,
		filepath.Join(base, constant.GelObjectsDirName),
		filepath.Join(base, constant.GelRefsDirName),
	}

	exists := initService.filesystemService.Exists(base)
	for _, dir := range dirs {
		if err := initService.filesystemService.MakeDir(dir, constant.GelDirectoryPermission); err != nil {
			return "", err
		}
	}

	if exists {
		return "Reinitialized existing Gel repository", nil
	}

	return "Initialized empty Gel repository in " + base, nil
}
