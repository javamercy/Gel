package core

import (
	"Gel/domain"
	"Gel/internal/workspace"
	"errors"
	"strings"
)

var (
	PathNotFoundInTreeError = errors.New("path not found")
)

type TreeResolver struct {
	objectService     *ObjectService
	indexService      *IndexService
	refService        *RefService
	pathResolver      *PathResolver
	hashObjectService *HashObjectService
}

func NewTreeResolver(
	objectService *ObjectService, indexService *IndexService, refService *RefService,
	pathResolver *PathResolver, hashObjectService *HashObjectService,
) *TreeResolver {
	return &TreeResolver{
		objectService:     objectService,
		indexService:      indexService,
		refService:        refService,
		pathResolver:      pathResolver,
		hashObjectService: hashObjectService,
	}
}

func (t *TreeResolver) ResolveHEAD() (map[string]string, error) {
	return t.ResolveRef(workspace.HeadFileName)
}

func (t *TreeResolver) ResolveRef(refName string) (map[string]string, error) {
	commitHash, err := t.refService.Resolve(refName)
	if err != nil {
		return nil, err
	}
	return t.ResolveCommit(commitHash)
}

func (t *TreeResolver) ResolveCommit(hash string) (map[string]string, error) {
	commit, err := t.objectService.ReadCommit(hash)
	if err != nil {
		return nil, err
	}

	entries := make(map[string]string)
	walker := NewTreeWalker(t.objectService, WalkOptions{Recursive: true})
	err = walker.Walk(
		commit.TreeHash, "", func(e domain.TreeEntry, relPath string) error {
			entries[relPath] = e.Hash
			return nil
		},
	)
	return entries, err
}

func (t *TreeResolver) ResolveIndex() (map[string]string, error) {
	entries, err := t.indexService.GetEntries()
	entriesMap := make(map[string]string)
	for _, entry := range entries {
		entriesMap[entry.Path] = entry.Hash
	}
	return entriesMap, err
}

func (t *TreeResolver) ResolveWorkingTree() (map[string]string, error) {
	resolvedPaths, err := t.pathResolver.Resolve([]string{"."})
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	for _, resolved := range resolvedPaths {
		for path := range resolved.NormalizedPaths {
			hash, _, err := t.hashObjectService.HashObject(path, false)
			if err != nil {
				return nil, err
			}
			results[path] = hash
		}
	}
	return results, nil
}

func (t *TreeResolver) LookupPathInTree(treeHash, path string) (domain.TreeEntry, error) {
	segments := strings.Split(path, "/")
	return t.lookupPathInTreeRecursive(treeHash, segments)
}

func (t *TreeResolver) lookupPathInTreeRecursive(treeHash string, segments []string) (domain.TreeEntry, error) {
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
	return domain.TreeEntry{}, PathNotFoundInTreeError
}
