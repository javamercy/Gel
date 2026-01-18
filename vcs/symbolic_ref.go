package vcs

import (
	"Gel/core/constant"
	"Gel/core/repository"
	"Gel/storage"
	"fmt"
	"path/filepath"
	"strings"
)

type SymbolicRefService struct {
	filesystemStorage  storage.IFilesystemStorage
	repositoryProvider repository.IRepositoryProvider
}

func NewSymbolicRefService(filesystemStorage storage.IFilesystemStorage, repositoryProvider repository.IRepositoryProvider) *SymbolicRefService {
	return &SymbolicRefService{
		filesystemStorage:  filesystemStorage,
		repositoryProvider: repositoryProvider,
	}
}

func (s *SymbolicRefService) Read(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("symbolic-ref: name is required")
	}

	repo := s.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, name)

	data, err := s.filesystemStorage.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(data))

	if !strings.HasPrefix(content, "ref: ") {
		return "", fmt.Errorf("%s is not a symbolic ref", name)
	}

	return strings.TrimPrefix(content, "ref: "), nil
}

func (s *SymbolicRefService) Update(name, ref string) error {
	if name == "" || ref == "" {
		return fmt.Errorf("symbolic-ref: name and ref are required")
	}

	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("refusing to point %s outside of refs/", name)
	}

	repo := s.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, name)

	content := fmt.Sprintf("ref: %s\n", ref)
	return s.filesystemStorage.WriteFile(path, []byte(content), false, constant.GelFilePermission)
}
