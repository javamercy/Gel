package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"Gel/storage"
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

var (
	pathNotFoundInTreeError = errors.New("path not found")
)

type RestoreService struct {
	indexService      *IndexService
	objectService     *ObjectService
	filesystemStorage *storage.FilesystemStorage
	refService        *RefService
}

func NewRestoreService(
	indexService *IndexService,
	objectService *ObjectService,
	filesystemStorage *storage.FilesystemStorage,
	refService *RefService) *RestoreService {
	return &RestoreService{
		indexService:      indexService,
		objectService:     objectService,
		filesystemStorage: filesystemStorage,
		refService:        refService,
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
	commitHash, err := r.refService.Resolve(constant.GelHeadFileName)
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
		existsInIndex, existsInHead := false, false

		for _, indexEntry := range index.Entries {
			if indexEntry.Path == path {
				existsInIndex = true
				break
			}
		}

		treeEntry, err := r.LookupPathInTree(commit.TreeHash, path)
		if err != nil && !errors.Is(err, pathNotFoundInTreeError) {
			return err
		} else if err == nil {
			existsInHead = true
		}

		if !existsInIndex && !existsInHead {
			return fmt.Errorf("pathspec %s not not match any files", path)
		}

		if existsInHead {
			// exists in Head, update or add to index
			newIndexEntry := domain.NewEmptyIndexEntry(path, treeEntry.Hash, treeEntry.Mode.Uint32())
			index.SetEntry(newIndexEntry)
		} else {
			// exists in index, remove it
			index.RemoveEntry(path)

		}
	}
	return r.indexService.Write(index)
}

func (r *RestoreService) restoreWorkingDir(paths []string) error {
	indexEntries, err := r.indexService.GetEntries()
	if err != nil {
		return err
	}

	for _, path := range paths {
		for _, entry := range indexEntries {
			if entry.Path == path {
				currHash, err := r.objectService.ComputeHash(path)
				if err != nil && !errors.Is(err, fs.ErrNotExist) {
					return err
				}
				if currHash != entry.Hash {
					// File in working dir has changed, restore it
					blob, err := r.objectService.ReadBlob(entry.Hash)
					if err != nil {
						return err
					}
					if err := r.filesystemStorage.WriteFile(
						path,
						blob.Body(),
						true,
						constant.GelFilePermission); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *RestoreService) LookupPathInTree(treeHash, path string) (domain.TreeEntry, error) {
	segments := strings.Split(path, "/")
	return r.lookupPathInTreeRecursive(treeHash, segments)
}

func (r *RestoreService) lookupPathInTreeRecursive(treeHash string, segments []string) (domain.TreeEntry, error) {
	entries, err := r.objectService.ReadTreeAndDeserializeEntries(treeHash)
	if err != nil {
		return domain.TreeEntry{}, err
	}
	for _, entry := range entries {
		if entry.Name == segments[0] {
			if len(segments) == 1 {
				return entry, nil
			}
			return r.lookupPathInTreeRecursive(entry.Hash, segments[1:])
		}
	}
	return domain.TreeEntry{}, pathNotFoundInTreeError
}
