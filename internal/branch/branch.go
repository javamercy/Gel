package branch

import (
	"Gel/internal/core"
	"Gel/internal/workspace"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type branch struct {
	name      string
	isCurrent bool
}
type BranchService struct {
	refService        *core.RefService
	objectService     *core.ObjectService
	workspaceProvider *workspace.Provider
}

func NewBranchService(
	refService *core.RefService, objectService *core.ObjectService, workspaceProvider *workspace.Provider,
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

	branches := make([]branch, 0)
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
			branches = append(
				branches, branch{
					name:      name,
					isCurrent: isCurrent,
				},
			)
			return nil
		},
	)
	if err != nil {
		return err
	}
	slices.SortFunc(
		branches, func(a, b branch) int {
			return strings.Compare(a.name, b.name)
		},
	)

	for _, b := range branches {
		if b.isCurrent {
			if _, err := fmt.Fprintf(writer, "%s* %s%s\n", core.ColorGreen, b.name, core.ColorReset); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(writer, "%s\n", b.name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BranchService) Create(name string, startPoint string) error {
	if err := validateBranchName(name); err != nil {
		return err
	}
	ref := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, name)
	exists := b.refService.Exists(ref)
	if exists {
		return errors.New("branch already exists")
	}

	if startPoint == "" {
		commitHash, err := b.refService.Resolve(workspace.HeadFileName)
		if err != nil {
			return err
		}
		return b.refService.Write(ref, commitHash)
	}

	if commitHash, err := b.refService.Resolve(
		filepath.Join(
			workspace.RefsDirName, workspace.HeadsDirName, startPoint,
		),
	); err == nil {
		return b.refService.Write(ref, commitHash)
	}

	_, err := b.objectService.ReadCommit(startPoint)
	if err != nil {
		return err
	}
	return b.refService.Write(ref, startPoint)
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
