package services

import (
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/encoding"
	"Gel/src/gel/core/utilities"
	"Gel/src/gel/domain"
	"Gel/src/gel/domain/objects"
	"Gel/src/gel/persistence/repositories"
	"bytes"
	"encoding/hex"
	"sort"
	"strings"
)

type IWriteTreeService interface {
	WriteTree() (string, *gelErrors.GelError)
}

type WriteTreeService struct {
	objectRepository repositories.IObjectRepository
	indexRepository  repositories.IIndexRepository
}

func NewWriteTreeService(objectRepository repositories.IObjectRepository, indexRepository repositories.IIndexRepository) *WriteTreeService {
	return &WriteTreeService{
		objectRepository: objectRepository,
		indexRepository:  indexRepository,
	}
}

func (writeTreeService *WriteTreeService) WriteTree() (string, *gelErrors.GelError) {
	entries, err := writeTreeService.indexRepository.GetEntries()
	if err != nil {
		return "", gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	root := buildTreeStructure(entries)
	rootHash, err := writeTreeService.buildTreeAndWrite(root)

	if err != nil {
		return "", gelErrors.NewGelError(gelErrors.ExitCodeFatal, err.Error())
	}

	return rootHash, nil
}

func (writeTreeService *WriteTreeService) buildTreeAndWrite(directory *DirectoryNode) (string, error) {
	var entries []*objects.TreeEntry

	for _, childDirectory := range directory.Children {
		subTreeHash, err := writeTreeService.buildTreeAndWrite(childDirectory)
		if err != nil {
			return "", err
		}
		entry := objects.NewTreeEntry(constant.GelDirectoryModeStr, subTreeHash, childDirectory.Name)
		entries = append(entries, entry)
	}

	for _, file := range directory.Files {
		entry := &objects.TreeEntry{
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

	treeObject := objects.NewTree(treeData)
	content := treeObject.Serialize()
	hash := encoding.ComputeHash(content)
	compressedContent, compressErr := encoding.Compress(content)
	if compressErr != nil {
		return "", compressErr
	}
	writeErr := writeTreeService.objectRepository.Write(hash, compressedContent)
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

func buildTreeData(entries []*objects.TreeEntry) ([]byte, error) {
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

func sortTreeEntries(entries []*objects.TreeEntry) {
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
