package vcs

import (
	"Gel/core/encoding"
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

func (writeTreeService *WriteTreeService) WriteTree() (string, error) {
	entries, err := writeTreeService.indexService.GetEntries()
	if err != nil {
		return "", err
	}

	root := buildRootTree(entries)
	rootHash, err := writeTreeService.buildTreeAndWrite(root)
	if err != nil {
		return "", err
	}
	return rootHash, nil
}

func (writeTreeService *WriteTreeService) buildTreeAndWrite(root *directoryNode) (string, error) {
	var entries []domain.TreeEntry

	for _, childDir := range root.children {
		subTreeHash, err := writeTreeService.buildTreeAndWrite(childDir)
		if err != nil {
			return "", err
		}
		entry, err := domain.NewTreeEntry(domain.Directory, subTreeHash, childDir.name)
		if err != nil {
			return "", err
		}
		entries = append(entries, entry)
	}

	for _, file := range root.files {
		entry, err := domain.NewTreeEntry(file.mode, file.hash, file.name)
		if err != nil {
			return "", err
		}
		entries = append(entries, entry)
	}

	sortTreeEntries(entries)

	tree, err := domain.NewTreeFromEntries(entries)
	if err != nil {
		return "", err
	}

	data := tree.Serialize()
	hash := encoding.ComputeSha256(data)
	err = writeTreeService.objectService.Write(hash, data)
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
		currDir := root
		names := strings.Split(entry.Path, "/")

		for i, name := range names {
			if i == len(names)-1 {
				fileNode := &fileNode{
					mode: domain.ParseFileMode(entry.Mode),
					hash: entry.Hash,
					name: name,
				}
				currDir.files = append(currDir.files, fileNode)
			} else {
				var childDir *directoryNode
				if existingChild, exists := currDir.children[name]; exists {
					childDir = existingChild
				} else {
					childDir = &directoryNode{
						name:     name,
						children: make(map[string]*directoryNode),
						files:    make([]*fileNode, 0),
					}
					currDir.children[name] = childDir
				}
				currDir = childDir
			}
		}
	}
	return root
}

func sortTreeEntries(entries []domain.TreeEntry) {
	sort.Slice(entries, func(i, j int) bool {
		NameI := entries[i].Name
		NameJ := entries[j].Name

		if entries[i].Mode.IsDirectory() {
			NameI += "/"
		}
		if entries[j].Mode.IsDirectory() {
			NameJ += "/"
		}
		return NameI < NameJ
	})
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
