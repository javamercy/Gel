package core

import (
	"Gel/domain"
	"Gel/internal/workspace"
	"strings"
)

type TreeResolver struct {
	objectService  *ObjectService
	indexService   *IndexService
	refService     *RefService
	pathResolver   *PathResolver
	changeDetector *ChangeDetector
}

func NewTreeResolver(
	objectService *ObjectService,
	indexService *IndexService,
	refService *RefService,
	pathResolver *PathResolver,
	changeDetector *ChangeDetector,
) *TreeResolver {
	return &TreeResolver{
		objectService:  objectService,
		indexService:   indexService,
		refService:     refService,
		pathResolver:   pathResolver,
		changeDetector: changeDetector,
	}
}

func (t *TreeResolver) ResolveHEAD() (map[string]domain.Hash, error) {
	return t.ResolveRef(workspace.HeadFileName)
}

func (t *TreeResolver) ResolveRef(refName string) (map[string]domain.Hash, error) {
	commitHash, err := t.refService.Resolve(refName)
	if err != nil {
		return nil, err
	}
	return t.ResolveCommit(commitHash)
}

func (t *TreeResolver) ResolveCommit(hash domain.Hash) (map[string]domain.Hash, error) {
	commit, err := t.objectService.ReadCommit(hash)
	if err != nil {
		return nil, err
	}

	entries := make(map[string]domain.Hash)
	walker := NewTreeWalker(t.objectService, WalkOptions{Recursive: true})
	err = walker.Walk(
		commit.TreeHash, "", func(e domain.TreeEntry, relPath string) error {
			entries[relPath] = e.Hash
			return nil
		},
	)
	return entries, err
}

func (t *TreeResolver) ResolveIndex() (map[string]domain.Hash, error) {
	entries, err := t.indexService.GetEntries()
	if err != nil {
		return nil, err
	}

	entriesMap := make(map[string]domain.Hash, len(entries))
	for _, entry := range entries {
		entriesMap[entry.Path.String()] = entry.Hash
	}
	return entriesMap, nil
}

func (t *TreeResolver) ResolveWorkingTree() (map[string]domain.Hash, error) {
	resolvedPaths, err := t.pathResolver.Resolve([]string{"."})
	if err != nil {
		return nil, err
	}

	index, err := t.indexService.Read()
	if err != nil {
		return nil, err
	}

	results := make(map[string]domain.Hash)
	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			fileStat := domain.GetFileStatFromPath(path.ToAbsolutePath())
			entry, _ := index.FindEntry(path.String())

			if entry != nil {
				changeResult, err := t.changeDetector.DetectFileChange(entry, fileStat)
				if err != nil {
					return nil, err
				}
				if !changeResult.IsModified {
					results[path.String()] = entry.Hash
				} else {
					results[path.String()] = changeResult.NewHash
				}
			} else {
				hash, _, err := t.objectService.ComputeObjectHash(path.ToAbsolutePath())
				if err != nil {
					return nil, err
				}
				results[path.String()] = hash
			}
		}
	}
	return results, nil
}

func (t *TreeResolver) LookupPathInTree(treeHash domain.Hash, path string) (domain.TreeEntry, error) {
	segments := strings.Split(path, "/")
	return t.lookupPathInTreeRecursive(treeHash, segments)
}

func (t *TreeResolver) lookupPathInTreeRecursive(treeHash domain.Hash, segments []string) (domain.TreeEntry, error) {
	entries, err := t.objectService.ReadTreeAndDeserializeEntries(treeHash)
	if err != nil {
		return domain.TreeEntry{}, err
	}
	for _, entry := range entries {
		if entry.Name == segments[0] {
			if len(segments) == 1 {
				return entry, nil
			}
			return t.lookupPathInTreeRecursive(entry.Hash, segments[1:])
		}
	}
	return domain.TreeEntry{}, ErrPathNotFoundInTree
}
