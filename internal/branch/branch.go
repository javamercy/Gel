package branch

import (
	"Gel/internal/core"
	"Gel/internal/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type BranchListItem struct {
	Name      string
	IsCurrent bool
}
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

func (b *BranchService) List() ([]BranchListItem, error) {
	headsDir := filepath.Join(b.workspace.GelDir, domain.RefsDirName, domain.HeadsDirName)
	currentBranchRef, err := b.refService.ReadSymbolic(domain.HeadFileName)
	if err != nil {
		return nil, fmt.Errorf("branch: failed to read symbolic ref: %w", err)
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
			isCurrent := ref == currentBranchRef
			branches[name] = isCurrent
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("branch: failed to list branches: %w", err)
	}

	branchNames := make([]BranchListItem, 0, len(branches))
	for name := range branches {
		branchNames = append(
			branchNames, BranchListItem{
				Name:      name,
				IsCurrent: branches[name],
			},
		)
	}
	slices.SortFunc(
		branchNames, func(a, b BranchListItem) int {
			return strings.Compare(a.Name, b.Name)
		},
	)
	return branchNames, nil
}

func (b *BranchService) Create(name string, startPoint string) error {
	if err := validateBranchName(name); err != nil {
		return fmt.Errorf("branch: '%s': %w", name, err)
	}
	if ok, err := b.Exists(name); err != nil {
		return err
	} else if ok {
		return fmt.Errorf("branch: '%s': %w", name, ErrBranchAlreadyExists)
	}

	ref := filepath.Join(domain.RefsDirName, domain.HeadsDirName, name)
	if startPoint == "" {
		commitHash, err := b.refService.Resolve(domain.HeadFileName)
		if err != nil {
			if errors.Is(err, core.ErrRefNotFound) {
				return fmt.Errorf("branch: '%s': %w", name, ErrNoCommitsYet)
			}
			return fmt.Errorf("branch: failed to resolve HEAD: %w", err)
		}
		if commitHash.IsEmpty() {
			return fmt.Errorf("branch: '%s': %w", name, ErrNoCommitsYet)
		}
		if err := b.refService.Write(ref, commitHash); err != nil {
			return fmt.Errorf("branch: failed to write '%s': %w", name, err)
		}
		return nil
	}

	startBranchRef := filepath.Join(domain.RefsDirName, domain.HeadsDirName, startPoint)
	if commitHash, err := b.refService.Read(startBranchRef); err == nil {
		return b.refService.Write(ref, commitHash)
	}

	startHash, err := domain.NewHash(startPoint)
	if err != nil {
		return fmt.Errorf("branch: '%s': %w", startPoint, ErrInvalidStartPoint)
	}
	if _, err := b.objectService.ReadCommit(startHash); err != nil {
		return fmt.Errorf("branch: '%s': %w", startPoint, ErrInvalidStartPoint)
	}
	if err := b.refService.Write(ref, startHash); err != nil {
		return fmt.Errorf("branch: %w", err)
	}
	return nil
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

func (b *BranchService) Exists(name string) (bool, error) {
	targetRef := filepath.Join(domain.RefsDirName, domain.HeadsDirName, name)
	ok, err := b.refService.Exists(targetRef)
	if err != nil {
		return false, fmt.Errorf("branch: %w", err)
	}
	return ok, nil
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
