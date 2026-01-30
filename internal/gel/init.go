package gel

import (
	"Gel/internal/workspace"
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

	base := filepath.Join(path, workspace.GelDirName)

	dirs := []string{
		filepath.Join(base, workspace.ObjectsDirName),
		filepath.Join(base, workspace.RefsDirName, workspace.HeadsDirName),
		filepath.Join(base, workspace.RefsDirName, workspace.TagsDirName),
	}

	files := []string{
		filepath.Join(base, workspace.ConfigFileName),
		filepath.Join(base, workspace.HeadFileName),
	}

	exists := initService.filesystemStorage.Exists(base)

	for _, dir := range dirs {
		if err := initService.filesystemStorage.MakeDir(
			dir,
			workspace.DirPermission); err != nil {
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
				workspace.FilePermission)
		} else {
			err = initService.filesystemStorage.WriteFile(
				file,
				[]byte{}, false,
				workspace.FilePermission)
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
