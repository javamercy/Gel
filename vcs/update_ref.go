package vcs

import (
	"Gel/core/constant"
	"Gel/core/repository"
	"Gel/storage"
	"Gel/vcs/validate"
	"fmt"
	"path/filepath"
)

type UpdateRefService struct {
	filesystemStorage  storage.IFilesystemStorage
	repositoryProvider repository.IRepositoryProvider
}

func NewUpdateRefService(filesystemStorage storage.IFilesystemStorage, repositoryProvider repository.IRepositoryProvider) *UpdateRefService {
	return &UpdateRefService{
		filesystemStorage:  filesystemStorage,
		repositoryProvider: repositoryProvider,
	}
}

func (s *UpdateRefService) Update(ref string, hash string) error {
	if err := validate.Hash(hash); err != nil {
		return err
	}

	repo := s.repositoryProvider.GetRepository()
	refPath := filepath.Join(repo.GelDir, ref)

	return s.filesystemStorage.WriteFile(
		refPath,
		[]byte(hash),
		true,
		constant.GelFilePermission)
}

func (s *UpdateRefService) SafeUpdate(ref string, newHash, oldHash string) error {
	if err := validate.Hash(newHash); err != nil {
		return err
	}
	if err := validate.Hash(oldHash); err != nil {
		return err
	}

	repo := s.repositoryProvider.GetRepository()
	refPath := filepath.Join(repo.GelDir, ref)
	data, err := s.filesystemStorage.ReadFile(refPath)
	if err != nil {
		return err
	}
	if string(data) != oldHash {
		return fmt.Errorf("ref %s is not pointing to %s", ref, oldHash)
	}

	return s.filesystemStorage.WriteFile(
		refPath,
		[]byte(newHash),
		true,
		constant.GelFilePermission)
}
