package branch

import (
	"Gel/internal/core"
	"Gel/internal/workspace"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type BranchService struct {
	refService        *core.RefService
	objectService     *core.ObjectService
	workspaceProvider *workspace.Provider
}

func NewBranchService(
	refService *core.RefService,
	objectService *core.ObjectService,
	workspaceProvider *workspace.Provider,
) *BranchService {
	return &BranchService{
		refService:        refService,
		objectService:     objectService,
		workspaceProvider: workspaceProvider,
	}
}

func (b *BranchService) List(writer io.Writer) error {
	ws := b.workspaceProvider.GetWorkspace()
	headsDir := filepath.Join(ws.GelDir, workspace.RefsDirName, workspace.HeadsDirName)

	currentBranch, err := b.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}

	branches := make(map[string]bool)
	err = filepath.WalkDir(
		headsDir, func(p string, d os.DirEntry, err error) error {
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
			branches[name] = isCurrent
			return nil
		},
	)
	if err != nil {
		return err
	}

	branchNames := make([]string, 0, len(branches))
	for name := range branches {
		branchNames = append(branchNames, name)
	}

	sort.Strings(branchNames)

	for _, name := range branchNames {
		if branches[name] {
			if _, err := fmt.Fprintf(writer, "%s* %s%s\n", core.ColorGreen, name, core.ColorReset); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(writer, "%s\n", name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BranchService) Create(branch string, startPoint string) error {
	if err := validateBranchName(branch); err != nil {
		return err
	}
	if b.Exists(branch) {
		return errors.New("branch already exists")
	}

	ref := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, branch)
	if startPoint == "" {
		commitHash, err := b.refService.Resolve(workspace.HeadFileName)
		if err != nil {
			return err
		}
		return b.refService.Write(ref, commitHash)
	}

	startBranchRef := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, startPoint)
	if commitHash, err := b.refService.Read(startBranchRef); err == nil {
		return b.refService.Write(ref, commitHash)
	}

	_, err := b.objectService.ReadCommit(startPoint)
	if err != nil {
		return err
	}
	return b.refService.Write(ref, startPoint)
}

func (b *BranchService) Delete(branch string) error {
	if err := validateBranchName(branch); err != nil {
		return err
	}

	currRef, err := b.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}

	refToDelete := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, branch)
	if refToDelete == currRef {
		return errors.New("cannot delete the current branch")
	}
	return b.refService.Delete(refToDelete)
}

func (b *BranchService) Exists(branch string) bool {
	targetRef := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, branch)
	return b.refService.Exists(targetRef)
}

func validateBranchName(branch string) error {
	if strings.HasPrefix(branch, "-") {
		return errors.New("branch name cannot start with '-'")
	}
	if strings.Contains(branch, "..") {
		return errors.New("branch name cannot contain '..'")
	}
	if strings.HasSuffix(branch, "/") {
		return errors.New("branch name cannot end with '/'")
	}
	return nil
}
