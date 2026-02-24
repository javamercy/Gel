package branch

import (
	"Gel/internal/core"
	"Gel/internal/inspect"
	"Gel/internal/tree"
	"Gel/internal/workspace"
	"fmt"
	"os"
	"path/filepath"
)

type SwitchService struct {
	indexService    *core.IndexService
	refService      *core.RefService
	branchService   *BranchService
	objectService   *core.ObjectService
	readTreeService *tree.ReadTreeService
	treeResolver    *core.TreeResolver
	restoreService  *inspect.RestoreService
}

func NewSwitchService(
	indexService *core.IndexService,
	refService *core.RefService,
	branchService *BranchService,
	objectService *core.ObjectService,
	readTreeService *tree.ReadTreeService,
	treeResolver *core.TreeResolver,
	restoreService *inspect.RestoreService,
) *SwitchService {
	return &SwitchService{
		indexService:    indexService,
		refService:      refService,
		branchService:   branchService,
		objectService:   objectService,
		readTreeService: readTreeService,
		treeResolver:    treeResolver,
		restoreService:  restoreService,
	}
}

func (s *SwitchService) Switch(branch string, create, force bool) (string, error) {
	if !force {
		if err := s.checkForUncommittedChanges(); err != nil {
			return "", err
		}
	}
	targetRef := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, branch)
	currentCommitHash, err := s.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return "", err
	}

	if create {
		if err := s.branchService.Create(branch, currentCommitHash); err != nil {
			return "", err
		}
	}
	if !s.branchService.Exists(branch) {
		return "", fmt.Errorf("branch '%s' does not exist", branch)
	}

	targetCommitHash, err := s.refService.Read(targetRef)
	if err != nil {
		return "", err
	}

	if currentCommitHash == targetCommitHash {
		return fmt.Sprintf("Switched to branch '%s'", branch),
			s.refService.WriteSymbolic(workspace.HeadFileName, targetRef)
	}

	if err := s.updateWorkingTree(currentCommitHash, targetCommitHash); err != nil {
		return "", err
	}

	targetCommit, err := s.objectService.ReadCommit(targetCommitHash)
	if err != nil {
		return "", err
	}
	if err := s.readTreeService.ReadTree(targetCommit.TreeHash); err != nil {
		return "", err
	}
	if err := s.refService.WriteSymbolic(workspace.HeadFileName, targetRef); err != nil {
		return "", err
	}

	headRef := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, workspace.HeadFileName)
	if err := s.refService.Write(headRef, targetCommitHash); err != nil {
		return "", err
	}
	return fmt.Sprintf("Switched to branch '%s'", branch), nil
}

func (s *SwitchService) updateWorkingTree(currentCommitHash, TargetCommitHash string) error {
	currentEntries, err := s.treeResolver.ResolveCommit(currentCommitHash)
	if err != nil {
		return err
	}
	targetEntries, err := s.treeResolver.ResolveCommit(TargetCommitHash)
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
	indexEntries, err := s.indexService.GetEntries()
	if err != nil {
		return err
	}
	headEntries, err := s.treeResolver.ResolveHEAD()
	if err != nil {
		return err
	}

	for _, indexEntry := range indexEntries {
		headHash, inHead := headEntries[indexEntry.Path]
		if !inHead || indexEntry.Hash != headHash {
			return fmt.Errorf("uncommitted changes in '%s'", indexEntry.Path)
		}
	}
	return nil
}
