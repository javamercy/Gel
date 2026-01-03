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

// Init initializes a new Gel repository at the specified path.
// It creates the necessary directory structure and configuration files.
// If a repository already exists at the path, it reinitializes it.
// It returns a message indicating the result of the operation.
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
