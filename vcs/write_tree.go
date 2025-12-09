package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
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
	var entries []*domain.TreeEntry

	for _, childDirectory := range directory.Children {
		subTreeHash, err := writeTreeService.buildTreeAndWrite(childDirectory)
		if err != nil {
			return "", err
		}
		entry := domain.NewTreeEntry(domain.Directory, subTreeHash, childDirectory.Name)
		entries = append(entries, entry)
	}

	for _, file := range directory.Files {
		entry := &domain.TreeEntry{
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

	treeObject := domain.NewTree(treeData)
	content := treeObject.Serialize()
	hash := encoding.ComputeHash(content)
	writeErr := writeTreeService.objectService.Write(hash, content)
	if writeErr != nil {
		return "", writeErr
	}
	return hash, nil
}
