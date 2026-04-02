package core

import (
	"Gel/internal/domain"
	"strings"
)

type TreeResolver struct {
	objectService  *ObjectService
	indexService   *IndexService
	refService     *RefService
	pathResolver   *PathResolver
	changeDetector *ChangeDetector
	workspace      *domain.Workspace
}

func NewTreeResolver(
	objectService *ObjectService,
	indexService *IndexService,
	refService *RefService,
	pathResolver *PathResolver,
	changeDetector *ChangeDetector,
	workspace *domain.Workspace,
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

func (t *TreeResolver) ResolveHEAD() (map[string]domain.Hash, error) {
	return t.ResolveRef(domain.HeadFileName)
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
			entry, _ := index.FindEntry(path)
			if entry != nil {
				changeResult, err := t.changeDetector.DetectFileChange(entry)
				if err != nil {
					return nil, err
				}
				switch changeResult.FileState {
				case FileStateUnchanged:
					results[path.String()] = entry.Hash
				case FileStateModified:
					results[path.String()] = changeResult.NewHash
				case FileStateDeleted:
					// TODO: what to do with deleted files?
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
