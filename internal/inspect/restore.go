package inspect

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/workspace"
	"errors"
	"os"
	"path/filepath"
)

type RestoreMode int

const (
	RestoreModeIndexVsWorkingTree RestoreMode = iota
	RestoreModeHEADVsIndex
	RestoreModeCommitVsWorkingTree
	RestoreModeCommitVsIndex
)

type RestoreOptions struct {
	Mode   RestoreMode
	Source string
}

type RestoreService struct {
	indexService   *core.IndexService
	objectService  *core.ObjectService
	refService     *core.RefService
	treeResolver   *core.TreeResolver
	changeDetector *core.ChangeDetector
}

func NewRestoreService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	refService *core.RefService,
	treeResolver *core.TreeResolver,
	changeDetector *core.ChangeDetector,
) *RestoreService {
	return &RestoreService{
		indexService:   indexService,
		objectService:  objectService,
		refService:     refService,
		treeResolver:   treeResolver,
		changeDetector: changeDetector,
	}
}

func (r *RestoreService) Restore(paths []string, options RestoreOptions) error {
	switch options.Mode {
	case RestoreModeIndexVsWorkingTree:
		return r.restoreIndexVsWorkingTree(paths)
	case RestoreModeHEADVsIndex:
		return r.restoreHEADVsIndex(paths)
	case RestoreModeCommitVsWorkingTree:
		commitHash, err := r.resolveSource(options.Source)
		if err != nil {
			return err
		}

		return r.restoreCommitVsWorkingTree(commitHash, paths)
	case RestoreModeCommitVsIndex:
		commitHash, err := r.resolveSource(options.Source)
		if err != nil {
			return err
		}
		return r.restoreCommitVsIndex(commitHash, paths)
	default:
		return errors.New("invalid restore mode")
	}
}

func (r *RestoreService) restoreIndexVsWorkingTree(paths []string) error {
	index, err := r.indexService.Read()
	if err != nil {
		return err
	}

	for _, path := range paths {
		stat := domain.GetFileStatFromPath(path)
		entry, _ := index.FindEntry(path)
		if entry == nil {
			continue
		}

		changeResult, err := r.changeDetector.DetectFileChange(entry, stat)
		if err != nil {
			return err
		}
		if !changeResult.IsModified {
			continue
		}

		blob, err := r.objectService.ReadBlob(entry.Hash)
		if err != nil {
			return err
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
			return err
		}
		if err := os.WriteFile(path, blob.Body(), workspace.FilePermission); err != nil {
			return err
		}
	}
	return nil
}

func (r *RestoreService) restoreHEADVsIndex(paths []string) error {
	commitHash, err := r.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return err
	}
	return r.restoreCommitVsIndex(commitHash, paths)
}

func (r *RestoreService) restoreCommitVsWorkingTree(commitHash string, paths []string) error {
	commitEntries, err := r.treeResolver.ResolveCommit(commitHash)
	if err != nil {
		return err
	}

	workingTreeEntries, err := r.treeResolver.ResolveWorkingTree()
	if err != nil {
		return err
	}

	for _, path := range paths {
		cHash, inCommit := commitEntries[path]
		if !inCommit {
			continue
		}

		workingTreeHash, inWorkingTree := workingTreeEntries[path]

		if !inWorkingTree || cHash != workingTreeHash {
			blob, err := r.objectService.ReadBlob(cHash)
			if err != nil {
				return err
			}
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, workspace.DirPermission); err != nil {
				return err
			}
			if err := os.WriteFile(path, blob.Body(), workspace.FilePermission); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RestoreService) restoreCommitVsIndex(commitHash string, paths []string) error {
	commit, err := r.objectService.ReadCommit(commitHash)
	if err != nil {
		return err
	}

	index, err := r.indexService.Read()
	if err != nil {
		return err
	}

	for _, path := range paths {
		treeEntry, err := r.treeResolver.LookupPathInTree(commit.TreeHash, path)
		if err != nil && !errors.Is(err, core.PathNotFoundInTreeError) {
			return err
		}

		inCommit := err == nil
		indexEntry, _ := index.FindEntry(path)
		inIndex := indexEntry != nil

		switch {
		case !inCommit && !inIndex:
			continue
		case inCommit && inIndex && indexEntry.Hash == treeEntry.Hash:
			continue
		case inCommit:
			newIndexEntry := domain.NewEmptyIndexEntry(path, treeEntry.Hash, treeEntry.Mode.Uint32())
			index.SetEntry(newIndexEntry)
		default:
			index.RemoveEntry(path)

		}
	}
	return r.indexService.Write(index)
}

func (r *RestoreService) resolveSource(source string) (string, error) {
	var commitHash string
	var err error

	switch source {
	case workspace.HeadFileName:
		commitHash, err = r.refService.Resolve(source)
	case workspace.MainBranchName:
		commitHash, err = r.refService.Read("refs/heads/main")
	default:
		commitHash = source
	}
	if err != nil {
		return "", err
	}
	return commitHash, nil
}
