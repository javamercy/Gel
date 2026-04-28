package inspect

import (
	"Gel/internal/branch"
	"Gel/internal/core"
	"Gel/internal/domain"
	"errors"
	"fmt"
)

// FileStatus describes one path and its status label.
type FileStatus struct {
	Path   domain.NormalizedPath
	Status string
}

// StatusResult contains categorized repository changes for status output.
type StatusResult struct {
	Staged        []FileStatus
	Unstaged      []FileStatus
	Untracked     []domain.NormalizedPath
	CurrentBranch string
	HeadTreeSize  int
}

// StatusService computes repository state from HEAD, index, and working tree snapshots.
type StatusService struct {
	indexService  *core.IndexService
	objectService *core.ObjectService
	branchService *branch.BranchService
	treeResolver  *core.TreeResolver
}

// NewStatusService creates a status service.
func NewStatusService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	branchService *branch.BranchService,
	treeResolver *core.TreeResolver,
) *StatusService {
	return &StatusService{
		indexService:  indexService,
		objectService: objectService,
		branchService: branchService,
		treeResolver:  treeResolver,
	}
}

// Status returns staged, unstaged, and untracked changes for the current branch.
func (s *StatusService) Status() (*StatusResult, error) {
	result := &StatusResult{}
	indexPathHashes, headTreePathHashes, workingTreePathHashes, err := s.loadPathHashes()
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}

	result.Staged = collectStaged(indexPathHashes, headTreePathHashes)
	result.Unstaged = collectUnstaged(indexPathHashes, workingTreePathHashes)
	result.Untracked = collectUntracked(indexPathHashes, workingTreePathHashes)

	currentBranch, err := s.branchService.Current()
	if err != nil {
		return nil, fmt.Errorf("status: %w", err)
	}
	result.CurrentBranch = currentBranch
	result.HeadTreeSize = len(headTreePathHashes)
	return result, nil
}

// loadPathHashes resolves path->hash snapshots from index, HEAD tree, and working tree.
func (s *StatusService) loadPathHashes() (
	indexPathHashes,
	headTreePathHashes,
	workingTreePathHashes core.PathHashes,
	err error,
) {
	indexPathHashes, err = s.treeResolver.ResolveIndex()
	if err != nil {
		return
	}

	headTreePathHashes, err = s.treeResolver.ResolveHEAD()
	if err != nil && !errors.Is(err, core.ErrRefNotFound) {
		return
	}

	workingTreePathHashes, err = s.treeResolver.ResolveWorkingTree()
	if err != nil {
		return
	}
	return
}

// collectStaged compares index against HEAD to find staged additions, modifications, and deletions.
func collectStaged(indexPathHashes, headTreePathHashes core.PathHashes) (staged []FileStatus) {
	for indexPath, indexHash := range indexPathHashes {
		headHash, inHead := headTreePathHashes[indexPath]
		if !inHead {
			staged = append(staged, FileStatus{indexPath, "New File"})
		} else if headHash != indexHash {
			staged = append(staged, FileStatus{indexPath, "Modified"})
		}
	}
	for path := range headTreePathHashes {
		if _, inIndex := indexPathHashes[path]; !inIndex {
			staged = append(staged, FileStatus{path, "Deleted"})
		}
	}
	return
}

// collectUnstaged compares working tree against index to find unstaged modifications and deletions.
func collectUnstaged(indexPathHashes, workingTreePathHashes core.PathHashes) (unstaged []FileStatus) {
	for indexPath, indexHash := range indexPathHashes {
		workingTreeHash, inWorkingDir := workingTreePathHashes[indexPath]
		if !inWorkingDir {
			// in Index but not in Working Dir
			unstaged = append(unstaged, FileStatus{indexPath, "Deleted"})
		} else if workingTreeHash != indexHash {
			// in Index and Working Dir but different
			unstaged = append(unstaged, FileStatus{indexPath, "Modified"})
		}
	}
	return
}

// collectUntracked finds working tree paths that are not present in index.
func collectUntracked(indexPathHashes, workingTreePathHashes core.PathHashes) (untracked []domain.NormalizedPath) {
	for path := range workingTreePathHashes {
		if _, inIndex := indexPathHashes[path]; !inIndex {
			untracked = append(untracked, path)
		}
	}
	return
}
