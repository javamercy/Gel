package core

import (
	domain2 "Gel/internal/domain"
	"strings"
)

type TreeResolver struct {
	objectService  *ObjectService
	indexService   *IndexService
	refService     *RefService
	pathResolver   *PathResolver
	changeDetector *ChangeDetector
	workspace      *domain2.Workspace
}

func NewTreeResolver(
	objectService *ObjectService,
	indexService *IndexService,
	refService *RefService,
	pathResolver *PathResolver,
	changeDetector *ChangeDetector,
	workspace *domain2.Workspace,
) *TreeResolver {
	return &TreeResolver{
		objectService:  objectService,
		indexService:   indexService,
		refService:     refService,
		pathResolver:   pathResolver,
		changeDetector: changeDetector,
		workspace:      workspace,
	}
}

func (t *TreeResolver) ResolveHEAD() (map[string]domain2.Hash, error) {
	return t.ResolveRef(domain2.HeadFileName)
}

func (t *TreeResolver) ResolveRef(refName string) (map[string]domain2.Hash, error) {
	commitHash, err := t.refService.Resolve(refName)
	if err != nil {
		return nil, err
	}
	return t.ResolveCommit(commitHash)
}

func (t *TreeResolver) ResolveCommit(hash domain2.Hash) (map[string]domain2.Hash, error) {
	commit, err := t.objectService.ReadCommit(hash)
	if err != nil {
		return nil, err
	}

	entries := make(map[string]domain2.Hash)
	walker := NewTreeWalker(t.objectService, WalkOptions{Recursive: true})
	err = walker.Walk(
		commit.TreeHash, "", func(e domain2.TreeEntry, relPath string) error {
			entries[relPath] = e.Hash
			return nil
		},
	)
	return entries, err
}

func (t *TreeResolver) ResolveIndex() (map[string]domain2.Hash, error) {
	entries, err := t.indexService.GetEntries()
	if err != nil {
		return nil, err
	}

	entriesMap := make(map[string]domain2.Hash, len(entries))
	for _, entry := range entries {
		entriesMap[entry.Path.String()] = entry.Hash
	}
	return entriesMap, nil
}

func (t *TreeResolver) ResolveWorkingTree() (map[string]domain2.Hash, error) {
	resolvedPaths, err := t.pathResolver.Resolve([]string{"."})
	if err != nil {
		return nil, err
	}

	index, err := t.indexService.Read()
	if err != nil {
		return nil, err
	}

	results := make(map[string]domain2.Hash)
	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			absolutePath, err := path.ToAbsolutePath(t.workspace.RepoDir)
			if err != nil {
				return nil, err
			}
			fileStat := domain2.GetFileStatFromPath(absolutePath)
			entry, _ := index.FindEntry(path)

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
				absolutePath, err := path.ToAbsolutePath(t.workspace.RepoDir)
				if err != nil {
					return nil, err
				}
				hash, _, err := t.objectService.ComputeObjectHash(absolutePath)
				if err != nil {
					return nil, err
				}
				results[path.String()] = hash
			}
		}
	}
	return results, nil
}

func (t *TreeResolver) LookupPathInTree(treeHash domain2.Hash, path string) (domain2.TreeEntry, error) {
	segments := strings.Split(path, "/")
	return t.lookupPathInTreeRecursive(treeHash, segments)
}

func (t *TreeResolver) lookupPathInTreeRecursive(treeHash domain2.Hash, segments []string) (domain2.TreeEntry, error) {
	entries, err := t.objectService.ReadTreeAndDeserializeEntries(treeHash)
	if err != nil {
		return domain2.TreeEntry{}, err
	}
	for _, entry := range entries {
		if entry.Name == segments[0] {
			if len(segments) == 1 {
				return entry, nil
			}
			return t.lookupPathInTreeRecursive(entry.Hash, segments[1:])
		}
	}
	return domain2.TreeEntry{}, ErrPathNotFoundInTree
}
