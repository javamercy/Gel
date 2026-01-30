package gel

import (
	"Gel/domain"
	"Gel/internal/gel/validate"
	"fmt"
	"io"
)

type CatFileService struct {
	objectService *ObjectService
}

func NewCatFileService(objectService *ObjectService) *CatFileService {
	return &CatFileService{
		objectService: objectService,
	}
}

func (c *CatFileService) CatFile(writer io.Writer, hash string, objectType, pretty, size, exists bool) error {
	if err := validate.Hash(hash); err != nil {
		return err
	}

	if exists {
		if !c.objectService.Exists(hash) {
			return fmt.Errorf("object %s does not exist", hash)
		}
		return nil
	}

	object, err := c.objectService.Read(hash)
	if err != nil {
		return err
	}

	if objectType {
		if _, err := fmt.Fprintf(writer, "%s\n", object.Type()); err != nil {
			return err
		}
	}
	if size {
		if _, err := fmt.Fprintf(writer, "%d\n", object.Size()); err != nil {
			return err
		}
	}
	if pretty {

		return catFileWithPretty(writer, object)
	}
	return nil
}

func catFileWithPretty(writer io.Writer, object domain.IObject) error {
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
			objectType, err := entry.Mode.ObjectType()
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(writer,
				"%s %s %s\t%s\n",
				entry.Mode,
				objectType,
				entry.Hash,
				entry.Name); err != nil {
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
		if _, err := fmt.Fprintf(writer,
			"%s %s\n",
			domain.CommitFieldTree,
			commit.TreeHash); err != nil {
			return err
		}
		for _, parentHash := range commit.ParentHashes {
			if _, err := fmt.Fprintf(writer,
				"%s %s\n",
				domain.CommitFieldParent,
				parentHash); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(writer,
			"%s %s <%s> %s %s\n"+
				"%s %s <%s> %s %s\n"+
				"\n%s\n",
			domain.CommitFieldAuthor,
			commit.Author.User.Name,
			commit.Author.User.Email,
			commit.Author.Timestamp,
			commit.Author.Timezone,
			domain.CommitFieldCommitter,
			commit.Committer.User.Name,
			commit.Committer.User.Email,
			commit.Committer.Timestamp,
			commit.Committer.Timezone,
			commit.Message,
		); err != nil {
			return err
		}
	}
	return nil
}
