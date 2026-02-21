package branch

import (
	"Gel/internal/core"
	"Gel/internal/tree"
	"Gel/internal/workspace"
	"fmt"
	"os"
	"path/filepath"
)

type SwitchService struct {
	indexService    *core.IndexService
	refService      *core.RefService
	objectService   *core.ObjectService
	readTreeService *tree.ReadTreeService
	treeResolver    *core.TreeResolver
}

func NewSwitchService(
	indexService *core.IndexService,
	refService *core.RefService,
	objectService *core.ObjectService,
	readTreeService *tree.ReadTreeService,
	treeResolver *core.TreeResolver,
) *SwitchService {
	return &SwitchService{
		indexService:    indexService,
		refService:      refService,
		objectService:   objectService,
		readTreeService: readTreeService,
		treeResolver:    treeResolver,
	}
}

func (s *SwitchService) Switch(branch string, create, force bool) (string, error) {
	if !force {
		if err := s.checkForUncommittedChanges(); err != nil {
			return "", err
		}
	}
	targetRef := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, branch)
	exists := s.refService.Exists(targetRef)
	currentCommitHash, err := s.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return "", err
	}

	if create {
		if exists {
			return "", fmt.Errorf("branch '%s' already exists", branch)
		}
		if err := s.refService.WriteSymbolic(workspace.HeadFileName, targetRef); err != nil {
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
	if err := s.updateWorkingTree(currentCommit.TreeHash, targetCommit.TreeHash); err != nil {
		return "", err
	}
	if err := s.readTreeService.ReadTree(targetCommit.TreeHash); err != nil {
		return "", err
	}
	if err := s.refService.WriteSymbolic(workspace.HeadFileName, targetRef); err != nil {
		return "", err
	}

	headRef := filepath.Join(workspace.RefsDirName, workspace.HeadFileName, workspace.HeadFileName)
	if err := s.refService.Write(headRef, targetCommitHash); err != nil {
		return "", err
	}
	return fmt.Sprintf("Switched to branch '%s'", branch), nil
}

func (s *SwitchService) updateWorkingTree(currentTreeHash, targetTreeHash string) error {
	currentEntries, err := s.treeResolver.ResolveCommit(currentTreeHash)
	if err != nil {
		return err
	}
	targetEntries, err := s.treeResolver.ResolveCommit(targetTreeHash)
	if err != nil {
		return err
	}
	for targetPath, targetHash := range targetEntries {
		currentHash, ok := currentEntries[targetPath]
		if !ok || targetHash != currentHash {
			blob, err := s.objectService.ReadBlob(targetHash)
			if err != nil {
				return err
			}
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
				return err
			}
			if err := os.WriteFile(targetPath, blob.Body(), workspace.FilePermission); err != nil {
				return err
			}
		}
	}

	for currentPath := range currentEntries {
		if _, existsInTarget := targetEntries[currentPath]; !existsInTarget {
			if err := os.RemoveAll(currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SwitchService) checkForUncommittedChanges() error {
	// TODO: we need a change detection service
	indexEntries, err := s.treeResolver.ResolveIndex()
	if err != nil {
		return err
	}
	headEntries, err := s.treeResolver.ResolveHEAD()
	if err != nil {
		return err
	}

	for indexPath, indexHash := range indexEntries {
		headHash, inHead := headEntries[indexPath]
		if !inHead || indexHash != headHash {
			return fmt.Errorf("uncommitted changes in '%s'", indexPath)
		}
	}
	return nil
}
