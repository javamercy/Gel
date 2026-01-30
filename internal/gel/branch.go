package gel

import (
	"Gel/internal/workspace"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type BranchService struct {
	refService        *RefService
	workspaceProvider *workspace.Provider
}

func NewBranchService(refService *RefService, workspaceProvider *workspace.Provider) *BranchService {
	return &BranchService{
		refService:        refService,
		workspaceProvider: workspaceProvider,
	}
}

func (b *BranchService) List() (map[string]bool, error) {
	ws := b.workspaceProvider.GetWorkspace()
	headsDir := filepath.Join(ws.GelDir, workspace.RefsDirName, workspace.HeadsDirName)

	currentBranch, err := b.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return nil, err
	}

	branchMap := make(map[string]bool)
	err = filepath.WalkDir(headsDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// TODO: Ensure the branch is valid

		ref := strings.TrimPrefix(p, ws.GelDir+"/")
		name := strings.TrimPrefix(p, headsDir+"/")
		isCurrent := ref == currentBranch
		branchMap[name] = isCurrent
		return nil
	})

	if err != nil {
		return nil, err
	}
	return branchMap, nil
}

func (b *BranchService) Create(name string) error {
	if err := validateBranchName(name); err != nil {
		return err
	}
	ref := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, name)
	exists := b.refService.Exists(ref)
	if exists {
		return errors.New("branch already exists")
	}

	commitHash, err := b.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return err
	}
	return b.refService.Write(ref, commitHash)
}

func (b *BranchService) Delete(name string) error {
	if err := validateBranchName(name); err != nil {
		return err
	}

	currRef, err := b.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}

	refToDelete := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, name)
	if refToDelete == currRef {
		return errors.New("cannot delete the current branch")
	}
	return b.refService.Delete(refToDelete)
}

func (b *BranchService) Exists(name string) error {
	if err := validateBranchName(name); err != nil {
		return err
	}
	ref := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, name)
	if exists := b.refService.Exists(ref); exists {
		return nil
	}
	return errors.New("branch does not exist")
}

func validateBranchName(name string) error {
	if strings.HasPrefix(name, "-") {
		return errors.New("branch name cannot start with '-'")
	}
	if strings.Contains(name, "..") {
		return errors.New("branch name cannot contain '..'")
	}
	if strings.HasSuffix(name, "/") {
		return errors.New("branch name cannot end with '/'")
	}
	return nil
}
