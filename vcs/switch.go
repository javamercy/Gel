package vcs

import (
	"Gel/core/constant"
	"Gel/core/repository"
	"Gel/domain"
	"Gel/storage"
	"fmt"
	"path/filepath"
)

type SwitchService struct {
	refService        *RefService
	objectService     *ObjectService
	filesystemStorage *storage.FilesystemStorage
	readTreeService   *ReadTreeService
	repoProvider      *repository.Provider
}

func NewSwitchService(
	refService *RefService,
	objectService *ObjectService,
	filesystemStorage *storage.FilesystemStorage,
	readTreeService *ReadTreeService,
	repoProvider *repository.Provider) *SwitchService {
	return &SwitchService{
		refService:        refService,
		objectService:     objectService,
		filesystemStorage: filesystemStorage,
		readTreeService:   readTreeService,
		repoProvider:      repoProvider,
	}
}

func (s *SwitchService) Switch(branch string, create, force bool) (string, error) {
	targetRef := filepath.Join(constant.GelRefsDirName, constant.GelHeadsDirName, branch)
	exists := s.refService.Exists(targetRef)
	currentCommitHash, err := s.refService.Resolve(constant.GelHeadFileName)
	if err != nil {
		return "", err
	}
	if create {
		if exists {
			return "", fmt.Errorf("branch '%s' already exists", branch)
		}
		if err := s.refService.WriteSymbolic(constant.GelHeadFileName, targetRef); err != nil {
			return "", err
		}
		if err := s.refService.Write(targetRef, currentCommitHash); err != nil {
			return "", err
		}
		return fmt.Sprintf("Created and switched to branch '%s'", branch), nil
	}
	if !exists {
		return "", fmt.Errorf("branch '%s' does not exist", branch)
	}

	currentCommit, err := s.objectService.ReadCommit(currentCommitHash)
	if err != nil {
		return "", err
	}

	targetCommitHash, err := s.refService.Read(targetRef)
	if err != nil {
		return "", err
	}

	targetCommit, err := s.objectService.ReadCommit(targetCommitHash)
	if err != nil {
		return "", err
	}
	if err := s.updateWorkingDir(currentCommit.TreeHash, targetCommit.TreeHash); err != nil {
		return "", err
	}
	if err := s.readTreeService.ReadTree(targetCommit.TreeHash); err != nil {
		return "", err
	}

	if err := s.refService.WriteSymbolic(constant.GelHeadFileName, targetRef); err != nil {
		return "", err
	}
	headRef := filepath.Join(constant.GelRefsDirName, constant.GelHeadFileName, constant.GelHeadFileName)
	if err := s.refService.Write(headRef, targetCommitHash); err != nil {
		return "", err
	}

	return fmt.Sprintf("Switched to branch '%s'", branch), nil
}

func (s *SwitchService) updateWorkingDir(currentTreeHash, targetTreeHash string) error {

	treeWalker := NewTreeWalker(s.objectService, WalkOptions{
		Recursive:    true,
		IncludeTrees: false,
		OnlyTrees:    false,
	})

	currentPathMap := make(map[string]string)
	currentProcessor := func(entry domain.TreeEntry, relPath string) error {
		currentPathMap[relPath] = entry.Hash
		return nil
	}
	err := treeWalker.Walk(currentTreeHash, "", currentProcessor)
	if err != nil {
		return err
	}

	targetPathMap := make(map[string]string)
	targetProcessor := func(entry domain.TreeEntry, relPath string) error {
		targetPathMap[relPath] = entry.Hash
		return nil
	}
	err = treeWalker.Walk(targetTreeHash, "", targetProcessor)
	if err != nil {
		return err
	}

	for targetPath, targetHash := range targetPathMap {
		currentHash, ok := currentPathMap[targetPath]
		if !ok || targetHash != currentHash {
			blob, err := s.objectService.ReadBlob(targetHash)
			if err != nil {
				return err
			}
			if err := s.filesystemStorage.WriteFile(
				targetPath,
				blob.Body(),
				true,
				constant.GelFilePermission); err != nil {
				return err
			}
		}
	}

	for currentPath := range currentPathMap {
		if _, existsInTarget := targetPathMap[currentPath]; !existsInTarget {
			if err := s.filesystemStorage.RemoveAll(currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}
