package rules

import (
	"Gel/persistence/repositories"
	"errors"
	"fmt"
	"os"
)

type HashObjectRules struct {
	filesystemRepository repositories.IFilesystemRepository
}

func NewHashObjectRules(filesystemRepository repositories.IFilesystemRepository) *HashObjectRules {
	return &HashObjectRules{
		filesystemRepository: filesystemRepository,
	}
}

func (hashObjectRules *HashObjectRules) PathsMustBeFiles(paths []string) error {
	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("'%s': cannot hash a directory", path))
		}
	}
	return nil
}

func (hashObjectRules *HashObjectRules) AllPathsMustExist(paths []string) error {
	for _, path := range paths {
		if !hashObjectRules.filesystemRepository.Exists(path) {
			return errors.New(fmt.Sprintf("cannot open '%s': no such file or directory", path))
		}
	}
	return nil
}
