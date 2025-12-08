package vcs

import (
	"Gel/core/constant"
	"Gel/core/utilities"
	"Gel/domain"
	"bytes"
	"encoding/hex"
	"path"
	"sort"
	"strings"
	"time"
)

type ReadTreeService struct {
	indexService  *IndexService
	objectService *ObjectService
}

func NewReadTreeService(indexService *IndexService, objectService *ObjectService) *ReadTreeService {
	return &ReadTreeService{
		indexService:  indexService,
		objectService: objectService,
	}
}

func (readTreeService *ReadTreeService) ReadTree(hash string) error {

	indexEntries, expandErr := readTreeService.expandTree(hash, "")
	if expandErr != nil {
		return expandErr
	}

	index := domain.NewEmptyIndex()
	for _, entry := range indexEntries {
		index.AddEntry(entry)
	}

	writeErr := readTreeService.indexService.Write(index)

	if writeErr != nil {
		return writeErr
	}

	return nil

}

func (readTreeService *ReadTreeService) expandTree(treeHash, prefix string) ([]*domain.IndexEntry, error) {

	result := make([]*domain.IndexEntry, 0)
	treeEntries, gelError := readTreeService.readTreeAndDeserializeTreeEntries(treeHash)
	if gelError != nil {
		return nil, gelError
	}

	for _, treeEntry := range treeEntries {
		objectType, err := domain.GetObjectTypeByMode(treeEntry.Mode)
		fullPath := path.Join(prefix, treeEntry.Name)
		if err != nil {
			return nil, err
		}

		if objectType == domain.GelTreeObjectType {
			indexEntries, err := readTreeService.expandTree(treeEntry.Hash, fullPath)
			if err != nil {
				return nil, err
			}
			result = append(result, indexEntries...)
		} else if objectType == domain.GelBlobObjectType {

			fileStatInfo, fileStatErr := utilities.GetFileStatFromPath(fullPath)
			if fileStatErr != nil {
				return nil, fileStatErr
			}

			size, sizeErr := readTreeService.objectService.GetObjectSize(treeEntry.Hash)
			if sizeErr != nil {
				return nil, sizeErr
			}

			indexEntry := domain.NewIndexEntry(
				fullPath,
				treeEntry.Hash,
				size,
				utilities.ConvertModeToUint32(treeEntry.Mode),
				fileStatInfo.Device,
				fileStatInfo.Inode,
				fileStatInfo.UserId,
				fileStatInfo.GroupId,
				domain.ComputeIndexFlags(fullPath, 0),
				time.Now(),
				time.Now())

			result = append(result, indexEntry)
		}

	}

	return result, nil
}

func (readTreeService *ReadTreeService) readTreeAndDeserializeTreeEntries(treeHash string) ([]*domain.TreeEntry, error) {
	object, err := readTreeService.objectService.Read(treeHash)
	if err != nil {
		return nil, err
	}

	tree, ok := object.(*domain.Tree)
	if !ok {
		return nil, err
	}

	treeEntries, err := tree.DeserializeTree()
	if err != nil {
		return nil, err
	}

	return treeEntries, nil
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

func (directoryNode *DirectoryNode) AddFile(file *FileNode) {
	directoryNode.Files = append(directoryNode.Files, file)
}

func (directoryNode *DirectoryNode) AddChildDirectory(child *DirectoryNode) {
	directoryNode.Children[child.Name] = child
}

func buildTreeStructure(entries []*domain.IndexEntry) *DirectoryNode {
	root := NewDirectoryNode("", map[string]*DirectoryNode{}, []*FileNode{})
	for _, entry := range entries {
		names := strings.Split(entry.Path, "/")

		currentDirectory := root
		for i, name := range names {
			if i == len(names)-1 {
				fileNode := NewFileNode(name, entry.Hash, utilities.ConvertModeToString(entry.Mode))
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

func buildTreeData(entries []*domain.TreeEntry) ([]byte, error) {
	var buffer bytes.Buffer
	for _, entry := range entries {
		buffer.WriteString(entry.Mode)
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

func sortTreeEntries(entries []*domain.TreeEntry) {
	sort.Slice(entries, func(i, j int) bool {
		NameI := entries[i].Name
		NameJ := entries[j].Name

		if entries[i].Mode == constant.GelDirectoryModeStr {
			NameI += constant.SlashStr
		}
		if entries[j].Mode == constant.GelDirectoryModeStr {
			NameJ += constant.SlashStr
		}
		return NameI < NameJ
	})
}
