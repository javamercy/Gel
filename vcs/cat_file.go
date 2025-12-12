package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"strconv"
	"strings"
)

type CatFileService struct {
	objectService *ObjectService
}

func NewCatFileService(objectService *ObjectService) *CatFileService {
	return &CatFileService{
		objectService: objectService,
	}
}

func (catFileService *CatFileService) CatFile(hash string, objectType, pretty, size, exists bool) (string, error) {
	object, err := catFileService.objectService.Read(hash)
	if err != nil || exists {
		return "", err
	}
	if objectType {
		return string(object.Type()), nil
	}

	if size {
		return strconv.Itoa(object.Size()), nil
	}

	if pretty {
		var result strings.Builder
		if object.IsTree() {
			tree, _ := object.(*domain.Tree)
			treeEntries, err := tree.DeserializeTree()
			if err != nil {
				return "", err
			}

			for _, entry := range treeEntries {
				result.WriteString(entry.Mode.String())
				result.WriteString(constant.SpaceStr)
				result.WriteString(string(object.Type()))
				result.WriteString(constant.SpaceStr)
				result.WriteString(entry.Hash)
				result.WriteString(constant.TabStr)
				result.WriteString(entry.Name)
				result.WriteString(constant.NewLineStr)
			}
		} else if object.IsBlob() {
			result.Write(object.Data())
		}
		return result.String(), nil
	}
	return "", err
}
