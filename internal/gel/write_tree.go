package gel

import (
	"Gel/domain"
	"sort"
	"strings"
)

type WriteTreeService struct {
	indexService  *IndexService
	objectService *ObjectService
}

func NewWriteTreeService(indexService *IndexService, objectService *ObjectService) *WriteTreeService {
	return &WriteTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (w *WriteTreeService) WriteTree() (string, error) {
	entries, err := w.indexService.GetEntries()
	if err != nil {
		return "", err
	}

	root := buildRootTree(entries)
	rootHash, err := w.writeTreeRecursive(root)
	if err != nil {
		return "", err
	}
	return rootHash, nil
}

func (w *WriteTreeService) writeTreeRecursive(root *directoryNode) (string, error) {
	var entries []domain.TreeEntry

	for _, childDir := range root.children {
		subTreeHash, err := w.writeTreeRecursive(childDir)
		if err != nil {
			return "", err
		}
		entry := domain.NewTreeEntry(domain.DirectoryMode, subTreeHash, childDir.name)
		entries = append(entries, entry)
	}

	for _, file := range root.files {
		entry := domain.NewTreeEntry(file.mode, file.hash, file.name)
		entries = append(entries, entry)
	}

	sortTreeEntries(entries)

	tree, err := domain.NewTreeFromEntries(entries)
	if err != nil {
		return "", err
	}

	data := tree.Serialize()
	hash := ComputeSHA256(data)
	if w.objectService.Exists(hash) {
		return hash, nil
	}

	err = w.objectService.Write(hash, data)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func buildRootTree(entries []*domain.IndexEntry) *directoryNode {
	root := &directoryNode{
		name:     "",
		children: make(map[string]*directoryNode),
		files:    make([]*fileNode, 0),
	}

	for _, entry := range entries {
		parentDir := root
		names := strings.Split(entry.Path, "/")

		for i, name := range names {
			if i == len(names)-1 {
				fileNode := &fileNode{
					mode: domain.ParseFileMode(entry.Mode),
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

func sortTreeEntries(entries []domain.TreeEntry) {
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
	mode domain.FileMode
	hash string
	name string
}
type directoryNode struct {
	name     string
	children map[string]*directoryNode
	files    []*fileNode
}
