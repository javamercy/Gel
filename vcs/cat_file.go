package vcs

import (
	"Gel/core/constant"
	"Gel/domain"
	"io"
	"strconv"
)

type CatFileService struct {
	objectService *ObjectService
}

func NewCatFileService(objectService *ObjectService) *CatFileService {
	return &CatFileService{
		objectService: objectService,
	}
}

func (catFileService *CatFileService) CatFile(writer io.Writer, hash string, objectType, pretty, size, exists bool) error {
	object, err := catFileService.objectService.Read(hash)
	if err != nil || exists {
		return err
	}

	if objectType {
		if _, err := io.WriteString(writer, string(object.Type())); err != nil {
			return err
		}
	}

	if size {
		if _, err := io.WriteString(writer, strconv.Itoa(object.Size())); err != nil {
			return err
		}
	}

	if pretty {
		switch object.Type() {
		case domain.ObjectTypeTree:
			tree, ok := object.(*domain.Tree)
			if !ok {
				return domain.ErrInvalidObjectType
			}
			treeEntries, err := tree.Deserialize()
			if err != nil {
				return err
			}

			for _, entry := range treeEntries {
				if _, err := io.WriteString(writer, entry.Mode.String()); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, string(tree.Type())); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, entry.Hash); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.TabStr); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, entry.Name); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
					return err
				}
			}
		case domain.ObjectTypeBlob:
			blob, ok := object.(*domain.Blob)
			if !ok {
				return domain.ErrInvalidObjectType
			}
			_, err := io.WriteString(writer, string(blob.Body()))
			return err
		case domain.ObjectTypeCommit:
			commit, ok := object.(*domain.Commit)
			if !ok {
				return domain.ErrInvalidObjectType
			}

			if _, err := io.WriteString(writer, domain.CommitFieldTree); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.TreeHash); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
				return err
			}

			for _, parentHash := range commit.ParentHashes {
				if _, err := io.WriteString(writer, domain.CommitFieldParent); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, parentHash); err != nil {
					return err
				}
				if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
					return err
				}
			}

			if _, err := io.WriteString(writer, domain.CommitFieldAuthor); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Author.User.Name); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.LessThanStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Author.User.Email); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.GreaterThanStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Author.Timestamp); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, domain.CommitFieldCommitter); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Committer.User.Name); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.LessThanStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Committer.User.Email); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.GreaterThanStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Committer.Timestamp); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.SpaceStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Committer.Timezone); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, constant.NewLineStr); err != nil {
				return err
			}
			if _, err := io.WriteString(writer, commit.Message); err != nil {
				return err
			}
		}
	}

	return nil
}
