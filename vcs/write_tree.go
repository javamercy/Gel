package vcs

import (
	"Gel/core/constant"
	"Gel/core/encoding"
	"Gel/domain"
	"bytes"
	"encoding/hex"
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

	root := buildTreeStructure(entries)
	rootHash, err := writeTreeService.buildTreeAndWrite(root)

	if err != nil {
		return "", err
	}
	return rootHash, nil
}

func (writeTreeService *WriteTreeService) buildTreeAndWrite(directory *DirectoryNode) (string, error) {
	var entries []domain.TreeEntry

	for _, childDirectory := range directory.Children {
		subTreeHash, err := writeTreeService.buildTreeAndWrite(childDirectory)
		if err != nil {
			return "", err
		}
		entry, err := domain.NewTreeEntry(domain.Directory, subTreeHash, childDirectory.Name)
		if err != nil {
			return "", err
		}
		entries = append(entries, entry)
	}

	for _, file := range directory.Files {
		entry := domain.TreeEntry{
			Mode: file.Mode,
			Hash: file.Hash,
			Name: file.Name,
		}
		entries = append(entries, entry)
	}

	sortTreeEntries(entries)

	treeData, buildTreeErr := buildTreeData(entries)
	if buildTreeErr != nil {
		return "", buildTreeErr
	}

	treeObject, err := domain.NewTree(treeData)
	if err != nil {
		return "", err
	}
	content := treeObject.Serialize()
	hash := encoding.ComputeHash(content)
	writeErr := writeTreeService.objectService.Write(hash, content)
	if writeErr != nil {
		return "", writeErr
	}
	return hash, nil
}

func buildTreeStructure(entries []*domain.IndexEntry) *DirectoryNode {
	root := NewDirectoryNode("", map[string]*DirectoryNode{}, []*FileNode{})
	for _, entry := range entries {
		names := strings.Split(entry.Path, "/")

		currentDirectory := root
		for i, name := range names {
			if i == len(names)-1 {
				fileNode := NewFileNode(domain.ParseFileMode(entry.Mode), entry.Hash, name)
				currentDirectory.AddFile(fileNode)
			} else {
				var childDirectory *DirectoryNode
				if existingChild, exists := currentDirectory.Children[name]; exists {
					childDirectory = existingChild
				} else {
					childDirectory = NewDirectoryNode(name, map[string]*DirectoryNode{}, []*FileNode{})
					currentDirectory.Children[name] = childDirectory
				}

				currentDirectory = childDirectory
			}
		}
	}
	return root
}

func buildTreeData(entries []domain.TreeEntry) ([]byte, error) {
	var buffer bytes.Buffer
	for _, entry := range entries {
		buffer.WriteString(entry.Mode.String())
		buffer.WriteString(constant.SpaceStr)
		buffer.WriteString(entry.Name)
		buffer.WriteString(constant.NullStr)

		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return nil, err
		}
		buffer.Write(hashBytes)
	}

	return buffer.Bytes(), nil
}

func sortTreeEntries(entries []domain.TreeEntry) {
	sort.Slice(entries, func(i, j int) bool {
		NameI := entries[i].Name
		NameJ := entries[j].Name

		if entries[i].Mode.IsDirectory() {
			NameI += constant.SlashStr
		}
		if entries[j].Mode.IsDirectory() {
			NameJ += constant.SlashStr
		}
		return NameI < NameJ
	})
}

type FileNode struct {
	Mode domain.FileMode
	Hash string
	Name string
}

func NewFileNode(mode domain.FileMode, hash, name string) *FileNode {
	{
		return &FileNode{
			Mode: mode,
			Hash: hash,
			Name: name,
		}
	}
}

type DirectoryNode struct {
	Name     string
	Children map[string]*DirectoryNode
	Files    []*FileNode
}

func NewDirectoryNode(name string, children map[string]*DirectoryNode, files []*FileNode) *DirectoryNode {
	return &DirectoryNode{
		Name:     name,
		Children: children,
		Files:    files,
	}
}

func (directoryNode *DirectoryNode) AddFile(file *FileNode) {
	directoryNode.Files = append(directoryNode.Files, file)
}

func (directoryNode *DirectoryNode) AddChildDirectory(child *DirectoryNode) {
	directoryNode.Children[child.Name] = child
}
