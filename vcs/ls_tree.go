package vcs

import (
	"Gel/core/constant"
	"errors"
	"path"
	"strings"
)

var notATreeError = errors.New("object is not a tree")

type LsTreeService struct {
	objectService *ObjectService
}

func NewLsTreeService(objectService *ObjectService) *LsTreeService {
	return &LsTreeService{
		objectService: objectService,
	}
}

func (lsTreeService *LsTreeService) LsTree(hash string, recursive, showTrees bool) (string, error) {

	// TODO: validate hash string

	var result strings.Builder
	if recursive {
		lsTreeEntries, err := lsTreeService.lsTreeWithRecursive(hash, "")
		if err != nil {
			return "", err
		}

		for _, entry := range lsTreeEntries {
			result.WriteString(entry.Mode)
			result.WriteString(constant.SpaceStr)
			result.WriteString(entry.ObjectType)
			result.WriteString(constant.SpaceStr)
			result.WriteString(entry.Hash)
			result.WriteString(constant.TabStr)
			result.WriteString(entry.FullPath)
			result.WriteString(constant.NewLineStr)
		}

		return result.String(), nil
	}
	entries, err := lsTreeService.objectService.ReadTreeAndDeserializeEntries(hash)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		objectType, err := entry.Mode.ObjectType()
		if err != nil {
			return "", err
		}
		result.WriteString(entry.Mode.String())
		result.WriteString(constant.SpaceStr)
		result.WriteString(string(objectType))
		result.WriteString(constant.SpaceStr)
		result.WriteString(entry.Hash)
		result.WriteString(constant.TabStr)
		result.WriteString(entry.Name)
		result.WriteString(constant.NewLineStr)
	}
	return result.String(), nil
}

func (lsTreeService *LsTreeService) lsTreeWithRecursive(treeHash, prefix string) ([]*LsTreeEntry, error) {
	result := make([]*LsTreeEntry, 0)
	treeEntries, err := lsTreeService.objectService.ReadTreeAndDeserializeEntries(treeHash)
	if err != nil {
		return nil, err
	}
	for _, entry := range treeEntries {
		fullPath := path.Join(prefix, entry.Name)
		if entry.Mode.IsDirectory() {
			lsTreeEntries, err := lsTreeService.lsTreeWithRecursive(entry.Hash, entry.Name)
			if err != nil {
				return nil, err
			}
			result = append(result, lsTreeEntries...)
		} else {
			objectType, _ := entry.Mode.ObjectType()
			lsTreeEntry := &LsTreeEntry{
				Mode:       entry.Mode.String(),
				ObjectType: string(objectType),
				Hash:       entry.Hash,
				FullPath:   fullPath,
			}
			result = append(result, lsTreeEntry)

		}

	}

	return result, nil
}

type LsTreeEntry struct {
	Mode       string
	ObjectType string
	Hash       string
	FullPath   string
}
