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
		switch object.Type() {
		case domain.ObjectTypeTree:
			tree, ok := object.(*domain.Tree)
			if !ok {
				return "", domain.ErrInvalidObjectType
			}
			treeEntries, err := tree.Deserialize()
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
		case domain.ObjectTypeBlob:
			blob, _ := object.(*domain.Blob)
			result.Write(blob.Body())
		case domain.ObjectTypeCommit:
			commit, ok := object.(*domain.Commit)
			if !ok {
				return "", domain.ErrInvalidObjectType
			}

			result.WriteString(domain.CommitFieldTree)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.TreeHash)
			result.WriteString(constant.NewLineStr)

			for _, parentHash := range commit.Fields.ParentHashes {
				result.WriteString(domain.CommitFieldParent)
				result.WriteString(constant.SpaceStr)
				result.WriteString(parentHash)
				result.WriteString(constant.NewLineStr)
			}

			result.WriteString(domain.CommitFieldAuthor)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Author.Name)
			result.WriteString(constant.SpaceStr)
			result.WriteString(constant.LessThanStr)
			result.WriteString(commit.Fields.Author.Email)
			result.WriteString(constant.GreaterThanStr)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Author.Timestamp)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Author.Timezone)
			result.WriteString(constant.NewLineStr)

			result.WriteString(domain.CommitFieldCommitter)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Committer.Name)
			result.WriteString(constant.SpaceStr)
			result.WriteString(constant.LessThanStr)
			result.WriteString(commit.Fields.Committer.Email)
			result.WriteString(constant.GreaterThanStr)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Committer.Timestamp)
			result.WriteString(constant.SpaceStr)
			result.WriteString(commit.Fields.Committer.Timezone)
			result.WriteString(constant.NewLineStr)
			result.WriteString(constant.NewLineStr)
			result.WriteString(commit.Fields.Message)

			return result.String(), nil
		}
	}
	return "", err
}
