package vcs

import (
	"Gel/core/constant"
	"Gel/core/repository"
	"Gel/storage"
	"Gel/vcs/validate"
	"fmt"
	"path/filepath"
	"strings"
)

type RefService struct {
	repositoryProvider repository.IRepositoryProvider
	filesystemStorage  storage.IFilesystemStorage
}

func NewRefService(repositoryProvider repository.IRepositoryProvider, filesystemStorage storage.IFilesystemStorage) *RefService {
	return &RefService{
		repositoryProvider: repositoryProvider,
		filesystemStorage:  filesystemStorage,
	}
}

func (s *RefService) ReadSymbolic(name string) (string, error) {
	repo := s.repositoryProvider.GetRepository()
	refPath := filepath.Join(repo.GelDir, name)
	contentBytes, err := s.filesystemStorage.ReadFile(refPath)
	if err != nil {
		return "", err
	}

	contentStr := strings.TrimSpace(string(contentBytes))
	if !strings.HasPrefix(contentStr, "ref: ") {
		return "", fmt.Errorf("invalid symbolic ref: %s", contentStr)
	}
	return strings.TrimPrefix(contentStr, "ref: "), nil
}

func (s *RefService) WriteSymbolic(name, ref string) error {
	if name == "" || ref == "" {
		return fmt.Errorf("symbolic-ref: name and ref are required")
	}

	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("symbolic-ref: ref must start with refs/")
	}

	repo := s.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, name)
	contentStr := fmt.Sprintf("ref: %s\n", ref)
	return s.filesystemStorage.WriteFile(path, []byte(contentStr), false, constant.GelFilePermission)
}

func (s *RefService) Read(ref string) (string, error) {
	if !strings.HasPrefix(ref, "refs/") {
		return "", fmt.Errorf("ref must start with refs/")
	}
	repo := s.repositoryProvider.GetRepository()
	absPath := filepath.Join(repo.GelDir, ref)
	contentBytes, err := s.filesystemStorage.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	hash := strings.TrimSpace(string(contentBytes))
	if err := validate.Hash(hash); err != nil {
		return "", err
	}
	return hash, nil
}
func (s *RefService) Write(ref, hash string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("ref must start with refs/")
	}
	if err := validate.Hash(hash); err != nil {
		return err
	}

	contentStr := fmt.Sprintf("%s\n", hash)
	repo := s.repositoryProvider.GetRepository()
	absPath := filepath.Join(repo.GelDir, ref)
	return s.filesystemStorage.WriteFile(absPath, []byte(contentStr), false, constant.GelFilePermission)
}

func (s *RefService) Resolve(name string) (string, error) {
	ref, err := s.ReadSymbolic(name)
	if err != nil {
		return "", err
	}
	return s.Read(ref)
}
