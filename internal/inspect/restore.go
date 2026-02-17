package inspect

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/workspace"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type RestoreService struct {
	indexService      *core.IndexService
	objectService     *core.ObjectService
	hashObjectService *core.HashObjectService
	refService        *core.RefService
	treeResolver      *core.TreeResolver
}

func NewRestoreService(
	indexService *core.IndexService,
	objectService *core.ObjectService,
	hashObjectService *core.HashObjectService,
	refService *core.RefService,
	treeResolver *core.TreeResolver,
) *RestoreService {
	return &RestoreService{
		indexService:      indexService,
		objectService:     objectService,
		hashObjectService: hashObjectService,
		refService:        refService,
		treeResolver:      treeResolver,
	}
}

func (r *RestoreService) Restore(paths []string, source string, staged bool) error {
	// TODO: handle source
	if staged {
		return r.restoreWithStaged(paths)
	}
	return r.restoreWorkingDir(paths)
}

func (r *RestoreService) restoreWithStaged(paths []string) error {
	commitHash, err := r.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return err
	}

	commit, err := r.objectService.ReadCommit(commitHash)
	if err != nil {
		return err
	}

	index, err := r.indexService.Read()
	if err != nil {
		return err
	}

	for _, path := range paths {
		existsInHead := false
		existsInIndex := index.HasEntry(path)

		treeEntry, err := r.treeResolver.LookupPathInTree(commit.TreeHash, path)
		if err != nil && !errors.Is(err, core.PathNotFoundInTreeError) {
			return err
		} else if err == nil {
			existsInHead = true
		}
		if !existsInIndex && !existsInHead {
			return errors.New("pathspec " + path + " not not match any files")
		}
		if existsInHead {
			newIndexEntry := domain.NewEmptyIndexEntry(path, treeEntry.Hash, treeEntry.Mode.Uint32())
			index.SetEntry(newIndexEntry)
		} else {
			index.RemoveEntry(path)
		}
	}
	return r.indexService.Write(index)
}

func (r *RestoreService) restoreWorkingDir(paths []string) error {
	index, err := r.indexService.Read()
	if err != nil {
		return err
	}

	for _, path := range paths {
		indexEntry, _ := index.FindEntry(path)
		if indexEntry == nil {
			continue
		}
		currHash, _, err := r.hashObjectService.HashObject(path, false)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if currHash != indexEntry.Hash {
			blob, err := r.objectService.ReadBlob(indexEntry.Hash)
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
