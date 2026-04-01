package tree

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
	"sort"
	"strings"
)

type WriteTreeService struct {
	indexService  *core.IndexService
	objectService *core.ObjectService
}

func NewWriteTreeService(indexService *core.IndexService, objectService *core.ObjectService) *WriteTreeService {
	return &WriteTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (w *WriteTreeService) WriteTree() (domain2.Hash, error) {
	entries, err := w.indexService.GetEntries()
	if err != nil {
		return domain2.Hash{}, err
	}

	root := buildRootTree(entries)
	rootHash, err := w.writeTreeRecursive(root)
	if err != nil {
		return domain2.Hash{}, err
	}
	return rootHash, nil
}

func (w *WriteTreeService) writeTreeRecursive(root *directoryNode) (domain2.Hash, error) {
	var entries []domain2.TreeEntry

	for _, childDir := range root.children {
		subTreeHash, err := w.writeTreeRecursive(childDir)
		if err != nil {
			return domain2.Hash{}, err
		}
		entry := domain2.NewTreeEntry(domain2.DirectoryMode, subTreeHash, childDir.name)
		entries = append(entries, entry)
	}

	for _, file := range root.files {
		entry := domain2.NewTreeEntry(file.mode, file.hash, file.name)
		entries = append(entries, entry)
	}

	sortTreeEntries(entries)

	tree, err := domain2.NewTreeFromEntries(entries)
	if err != nil {
		return domain2.Hash{}, err
	}

	data := tree.Serialize()
	hexHash := core.ComputeSHA256(data)
	hash, err := domain2.NewHash(hexHash)
	if err != nil {
		return domain2.Hash{}, err
	}

	ok, err := w.objectService.Exists(hash)
	if err != nil {
		return domain2.Hash{}, err
	}
	if ok {
		return hash, nil
	}

	err = w.objectService.Write(hash, data)
	if err != nil {
		return domain2.Hash{}, err
	}
	return hash, nil
}

func buildRootTree(entries []*domain2.IndexEntry) *directoryNode {
	root := &directoryNode{
		name:     "",
		children: make(map[string]*directoryNode),
		files:    make([]*fileNode, 0),
	}

	for _, entry := range entries {
		parentDir := root
		names := strings.Split(entry.Path.String(), "/")

		for i, name := range names {
			if i == len(names)-1 {
				fileNode := &fileNode{
					mode: domain2.ParseFileMode(entry.Mode),
					hash: entry.Hash,
					name: name,
				}
				parentDir.files = append(parentDir.files, fileNode)
			} else {
				var childDir *directoryNode
				if existingChild, exists := parentDir.children[name]; exists {
					childDir = existingChild
				} else {
					childDir = &directoryNode{
						name:     name,
						children: make(map[string]*directoryNode),
						files:    make([]*fileNode, 0),
					}
					parentDir.children[name] = childDir
				}
				parentDir = childDir
			}
		}
	}
	return root
}

func sortTreeEntries(entries []domain2.TreeEntry) {
	sort.Slice(
		entries, func(i, j int) bool {
			NameI := entries[i].Name
			NameJ := entries[j].Name

			if entries[i].Mode.IsDirectory() {
				NameI += "/"
			}
			if entries[j].Mode.IsDirectory() {
				NameJ += "/"
			}
			return NameI < NameJ
		},
	)
}

type fileNode struct {
	mode domain2.FileMode
	hash domain2.Hash
	name string
}
type directoryNode struct {
	name     string
	children map[string]*directoryNode
	files    []*fileNode
}
