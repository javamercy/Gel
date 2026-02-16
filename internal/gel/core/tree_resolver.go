package core

import (
	"Gel/domain"
	"Gel/internal/gel/workspace"
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
