package internal

import (
	"Gel/internal/core"
	"Gel/internal/domain"
	"Gel/internal/tree"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ResetMode int

const (
	ResetModeSoft ResetMode = iota
	ResetModeMixed
	ResetModeHard
)

type ResetOptions struct {
	Mode ResetMode
}

type ResetResult struct {
	TargetHash domain.Hash
}

type ResetService struct {
	refService      *core.RefService
	objectService   *core.ObjectService
	readTreeService *tree.ReadTreeService
	treeResolver    *core.TreeResolver
	commitResolver  *core.CommitResolver
	workspace       *domain.Workspace
}

func NewResetService(
	refService *core.RefService,
	objectService *core.ObjectService,
	readTreeService *tree.ReadTreeService,
	treeResolver *core.TreeResolver,
	commitResolver *core.CommitResolver,
	workspace *domain.Workspace,
) *ResetService {
	return &ResetService{
		refService:      refService,
		objectService:   objectService,
		readTreeService: readTreeService,
		treeResolver:    treeResolver,
		commitResolver:  commitResolver,
		workspace:       workspace,
	}
}

func (r *ResetService) Reset(target string, options ResetOptions) (*ResetResult, error) {
	if err := validateMode(options.Mode); err != nil {
		return nil, fmt.Errorf("reset: %w", err)
	}

	if strings.TrimSpace(target) == "" {
		target = domain.HeadFileName
	}

	targetHash, err := r.commitResolver.Resolve(target)
	if err != nil {
		return nil, fmt.Errorf("reset: %w", err)
	}

	targetCommit, err := r.objectService.ReadCommit(targetHash)
	if err != nil {
		return nil, fmt.Errorf("reset: %w", err)
	}
	if options.Mode == ResetModeMixed || options.Mode == ResetModeHard {
		if err := r.readTreeService.ReadTree(targetCommit.TreeHash); err != nil {
			return nil, fmt.Errorf("reset: %w", err)
		}
	}
	if options.Mode == ResetModeHard {
		if err := r.checkoutWorkingTree(targetHash); err != nil {
			return nil, fmt.Errorf("reset: %w", err)
		}
	}
	if err := r.moveHEADPointer(targetHash); err != nil {
		return nil, fmt.Errorf("reset: %w", err)
	}
	return &ResetResult{
		TargetHash: targetHash,
	}, nil
}

func (r *ResetService) moveHEADPointer(hash domain.Hash) error {
	ref, err := r.refService.ReadSymbolic(domain.HeadFileName)
	if err != nil {
		return err
	}
	return r.refService.Write(ref, hash)
}

func (r *ResetService) checkoutWorkingTree(targetHash domain.Hash) error {
	headPathHashes, err := r.treeResolver.ResolveHEAD()
	if err != nil {
		return err
	}

	targetPathHashes, err := r.treeResolver.ResolveCommit(targetHash)
	if err != nil {
		return err
	}

	workingTreePathHashes, err := r.treeResolver.ResolveWorkingTree()
	if err != nil {
		return err
	}

	deletePaths := collectDeletePaths(headPathHashes, targetPathHashes)
	writePaths := collectWritePaths(targetPathHashes, workingTreePathHashes)

	for _, path := range deletePaths {
		absPath, err := path.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return err
		}
		if err := os.Remove(absPath.String()); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to delete '%s': %w", absPath, err)
		}
		if err := pruneEmptyParentDirs(absPath.String(), r.workspace.RepoDir); err != nil {
			return fmt.Errorf("failed to prune empty parents for '%s': %w", absPath, err)
		}
	}
	for _, path := range writePaths {
		blobHash := targetPathHashes[path]
		blob, err := r.objectService.ReadBlob(blobHash)
		if err != nil {
			return err
		}

		absPath, err := path.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(absPath.String()), domain.DirPermission); err != nil {
			return fmt.Errorf("failed to create parent dir for '%s': %w", absPath.String(), err)
		}
		if err := os.WriteFile(absPath.String(), blob.Body(), domain.FilePermission); err != nil {
			return fmt.Errorf("failed to write '%s': %w", absPath, err)
		}
	}
	return nil
}

// pruneEmptyParentDirs removes empty ancestor directories up to repoRoot.
// This keeps hard reset working when deleting the last file from a directory
// and then recreating that path as a regular file.
func pruneEmptyParentDirs(filePath, repoRoot string) error {
	dir := filepath.Dir(filePath)
	repoRoot = filepath.Clean(repoRoot)

	for dir != repoRoot {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		if len(entries) != 0 {
			return nil
		}

		if err := os.Remove(dir); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		dir = filepath.Dir(dir)
	}

	return nil
}

func collectDeletePaths(oldPathHashes, newPathHashes core.PathHashes) (deletePaths []domain.NormalizedPath) {
	for path := range oldPathHashes {
		if _, inNew := newPathHashes[path]; !inNew {
			deletePaths = append(deletePaths, path)
		}
	}
	return
}

func collectWritePaths(targetPathHashes, workingTreePathHashes core.PathHashes) (writePaths []domain.NormalizedPath) {
	for path, targetHash := range targetPathHashes {
		if workingHash, inWorking := workingTreePathHashes[path]; !inWorking || workingHash != targetHash {
			writePaths = append(writePaths, path)
		}
	}
	return
}

func validateMode(mode ResetMode) error {
	switch mode {
	case ResetModeSoft, ResetModeMixed, ResetModeHard:
		return nil
	default:
		return fmt.Errorf("invalid reset mode: %d", mode)
	}
}
