package branch

import (
	"Gel/domain"
	core2 "Gel/internal/core"
	"Gel/internal/tree"
	"Gel/internal/workspace"
	"fmt"
	"os"
	"path/filepath"
)

type SwitchService struct {
	refService        *core2.RefService
	objectService     *core2.ObjectService
	readTreeService   *tree.ReadTreeService
	workspaceProvider *workspace.Provider
}

func NewSwitchService(
	refService *core2.RefService,
	objectService *core2.ObjectService,
	readTreeService *tree.ReadTreeService,
	workspaceProvider *workspace.Provider,
) *SwitchService {
	return &SwitchService{
		refService:        refService,
		objectService:     objectService,
		readTreeService:   readTreeService,
		workspaceProvider: workspaceProvider,
	}
}

func (s *SwitchService) Switch(branch string, create, force bool) (string, error) {
	// TODO: handle force, also there are some improvements need to be done.
	// I'll get back here after implementing Status/Diff commands.

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
	if err := s.updateWorkingDir(currentCommit.TreeHash, targetCommit.TreeHash); err != nil {
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

func (s *SwitchService) updateWorkingDir(currentTreeHash, targetTreeHash string) error {
	treeWalker := core2.NewTreeWalker(
		s.objectService, core2.WalkOptions{
			Recursive:    true,
			IncludeTrees: false,
			OnlyTrees:    false,
		},
	)
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
			dir := filepath.Dir(targetPath)
			if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
				return err
			}
			if err := os.WriteFile(targetPath, blob.Body(), workspace.FilePermission); err != nil {
				return err
			}
		}
	}

	for currentPath := range currentPathMap {
		if _, existsInTarget := targetPathMap[currentPath]; !existsInTarget {
			if err := os.RemoveAll(currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}
