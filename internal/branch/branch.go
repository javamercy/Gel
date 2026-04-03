package branch

import (
	"Gel/internal/core"
	"Gel/internal/domain"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type BranchService struct {
	refService    *core.RefService
	objectService *core.ObjectService
	workspace     *domain.Workspace
}

func NewBranchService(
	refService *core.RefService,
	objectService *core.ObjectService,
	workspace *domain.Workspace,
) *BranchService {
	return &BranchService{
		refService:    refService,
		objectService: objectService,
		workspace:     workspace,
	}
}

func (b *BranchService) List(writer io.Writer) error {
	headsDir := filepath.Join(b.workspace.GelDir, domain.RefsDirName, domain.HeadsDirName)

	currentBranch, err := b.refService.ReadSymbolic(domain.HeadFileName)
	if err != nil {
		return fmt.Errorf("branch: failed to read symbolic ref: %w", err)
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

			ref := strings.TrimPrefix(p, b.workspace.GelDir+"/")
			name := strings.TrimPrefix(p, headsDir+"/")
			isCurrent := ref == currentBranch
			branches[name] = isCurrent
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("branch: failed to list branches: %w", err)
	}

	branchNames := make([]string, 0, len(branches))
	for name := range branches {
		branchNames = append(branchNames, name)
	}

	sort.Strings(branchNames)

	for _, name := range branchNames {
		if branches[name] {
			if _, err := fmt.Fprintf(writer, "%s* %s%s\n", core.ColorGreen, name, core.ColorReset); err != nil {
				return fmt.Errorf("branch: failed to write branch name: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(writer, "%s\n", name); err != nil {
				return fmt.Errorf("branch: failed to write branch name: %w", err)
			}
		}
	}
	return nil
}

func (b *BranchService) Create(name string, startPoint string) error {
	if err := validateBranchName(name); err != nil {
		return fmt.Errorf("'%s': %w", name, err)
	}

	if b.Exists(name) {
		return fmt.Errorf("'%s': %w", name, ErrBranchAlreadyExists)
	}

	ref := filepath.Join(domain.RefsDirName, domain.HeadsDirName, name)
	if startPoint == "" {
		commitHash, err := b.refService.Resolve(domain.HeadFileName)
		if err != nil {
			return err
		}
		return b.refService.Write(ref, commitHash)
	}

	startBranchRef := filepath.Join(domain.RefsDirName, domain.HeadsDirName, startPoint)
	if commitHash, err := b.refService.Read(startBranchRef); err == nil {
		return b.refService.Write(ref, commitHash)
	}

	startHash, err := domain.NewHash(startPoint)
	if err != nil {
		return err
	}
	if _, err := b.objectService.ReadCommit(startHash); err != nil {
		return err
	}
	return b.refService.Write(ref, startHash)
}

func (b *BranchService) Delete(name string) error {
	if err := validateBranchName(name); err != nil {
		return fmt.Errorf("'%s': %w", name, err)
	}

	currRef, err := b.refService.ReadSymbolic(domain.HeadFileName)
	if err != nil {
		return err
	}

	refToDelete := filepath.Join(domain.RefsDirName, domain.HeadsDirName, name)
	if refToDelete == currRef {
		return fmt.Errorf("'%s': %w", name, ErrDeleteCurrentBranch)
	}
	return b.refService.Delete(refToDelete)
}

func (b *BranchService) Exists(name string) bool {
	targetRef := filepath.Join(domain.RefsDirName, domain.HeadsDirName, name)
	ok, err := b.refService.Exists(targetRef)
	return ok && err == nil
}

func validateBranchName(name string) error {
	switch {
	case strings.HasPrefix(name, "-"):
		return fmt.Errorf("must not start with '-': %w", ErrInvalidBranchName)
	case strings.Contains(name, ".."):
		return fmt.Errorf("must not contain '..': %w", ErrInvalidBranchName)
	case strings.HasSuffix(name, "/"):
		return fmt.Errorf("must not end with '/': %w", ErrInvalidBranchName)
	}
	return nil
}
