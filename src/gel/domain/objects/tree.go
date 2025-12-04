package objects

import (
	"Gel/src/gel/core/constant"
	"bytes"
	"encoding/hex"
	"sort"
)

type Tree struct {
	*BaseObject
}

type TreeEntry struct {
	Mode string
	Hash string
	Name string
}

func NewTree(entries []*TreeEntry) *Tree {
	sortTreeEntries(entries)
	data := buildTreeData(entries)
	return &Tree{
		&BaseObject{
			objectType: GelTreeObjectType,
			data:       data,
		},
	}
}

func buildTreeData(entries []*TreeEntry) []byte {
	var buffer bytes.Buffer
	for _, entry := range entries {
		buffer.WriteString(entry.Mode)
		buffer.WriteByte(constant.SpaceByte)
		buffer.WriteString(entry.Name)
		buffer.WriteByte(constant.NullByte)

		hashBytes, _ := hex.DecodeString(entry.Hash)
		buffer.Write(hashBytes)
	}
	return buffer.Bytes()
}

func sortTreeEntries(entries []*TreeEntry) {
	sort.Slice(entries, func(i, j int) bool {
		NameI := entries[i].Name
		NameJ := entries[j].Name

		if entries[i].Mode == constant.GelDirMode {
			NameI += constant.SlashStr
		}
		if entries[j].Mode == constant.GelDirMode {
			NameJ += constant.SlashStr
		}
		return NameI < NameJ
	})
}

type TreeBuilder struct {
}

func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{}
}

type FileNode struct {
	Name string
	Hash string
	Mode string
}

func NewFileNode(name, hash, mode string) *FileNode {
	{
		return &FileNode{
			Name: name,
			Hash: hash,
			Mode: mode,
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

func (directory *DirectoryNode) AddFile(file *FileNode) {
	directory.Files = append(directory.Files, file)
}

func (directory *DirectoryNode) AddChildDirectory(child *DirectoryNode) {
	directory.Children[child.Name] = child
}
