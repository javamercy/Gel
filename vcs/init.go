package vcs

import (
	"Gel/core/constant"
	"fmt"
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

	base := filepath.Join(path, constant.GelRepositoryName)

	dirs := []string{
		base,
		filepath.Join(base, constant.GelObjectsDirectoryName),
		filepath.Join(base, constant.GelRefsDirectoryName),
	}

	exists := initService.filesystemService.Exists(base)
	for _, dir := range dirs {
		if err := initService.filesystemService.MakeDirectory(dir, constant.GelDirectoryPermission); err != nil {
			return "", err
		}
	}

	if exists {
		return "Reinitialized existing Gel repository", nil
	}

	return fmt.Sprintf("Initialized empty Gel repository in %v", base), nil
}
