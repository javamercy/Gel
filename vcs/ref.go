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
	repositoryProvider *repository.Provider
	filesystemStorage  *storage.FilesystemStorage
}

func NewRefService(repositoryProvider *repository.Provider, filesystemStorage *storage.FilesystemStorage) *RefService {
	return &RefService{
		repositoryProvider: repositoryProvider,
		filesystemStorage:  filesystemStorage,
	}
}

func (r *RefService) ReadSymbolic(name string) (string, error) {
	repo := r.repositoryProvider.GetRepository()
	refPath := filepath.Join(repo.GelDir, name)
	contentBytes, err := r.filesystemStorage.ReadFile(refPath)
	if err != nil {
		return "", err
	}

	contentStr := strings.TrimSpace(string(contentBytes))
	if !strings.HasPrefix(contentStr, "ref: ") {
		return "", fmt.Errorf("invalid symbolic ref: %s", contentStr)
	}
	return strings.TrimPrefix(contentStr, "ref: "), nil
}

func (r *RefService) WriteSymbolic(name, ref string) error {
	if name == "" || ref == "" {
		return fmt.Errorf("symbolic-ref: name and ref are required")
	}

	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("symbolic-ref: ref must start with refs/")
	}

	repo := r.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, name)
	contentStr := fmt.Sprintf("ref: %s\n", ref)
	return r.filesystemStorage.WriteFile(path, []byte(contentStr), false, constant.GelFilePermission)
}

func (r *RefService) Read(ref string) (string, error) {
	if !strings.HasPrefix(ref, "refs/") {
		return "", fmt.Errorf("ref must start with refs/")
	}
	repo := r.repositoryProvider.GetRepository()
	absPath := filepath.Join(repo.GelDir, ref)
	contentBytes, err := r.filesystemStorage.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	hash := strings.TrimSpace(string(contentBytes))
	if err := validate.Hash(hash); err != nil {
		return "", err
	}
	return hash, nil
}
func (r *RefService) Write(ref, hash string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("ref must start with refs/")
	}
	if err := validate.Hash(hash); err != nil {
		return err
	}

	contentStr := fmt.Sprintf("%s\n", hash)
	repo := r.repositoryProvider.GetRepository()
	absPath := filepath.Join(repo.GelDir, ref)
	return r.filesystemStorage.WriteFile(absPath, []byte(contentStr), true, constant.GelFilePermission)
}

func (r *RefService) Delete(ref string) error {
	if !strings.HasPrefix(ref, "refs/") {
		return fmt.Errorf("ref must start with refs/")
	}

	repo := r.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, ref)
	return r.filesystemStorage.RemoveAll(path)
}

func (r *RefService) Exists(ref string) bool {
	repo := r.repositoryProvider.GetRepository()
	path := filepath.Join(repo.GelDir, ref)
	return r.filesystemStorage.Exists(path)
}

func (r *RefService) Resolve(name string) (string, error) {
	ref, err := r.ReadSymbolic(name)
	if err != nil {
		return "", err
	}
	return r.Read(ref)
}
