package gel

import (
	"Gel/internal/workspace"
	"errors"
)

type CommitService struct {
	writeTreeService  *WriteTreeService
	commitTreeService *CommitTreeService
	refService        *RefService
	objectService     *ObjectService
}

func NewCommitService(
	writeTreeService *WriteTreeService,
	commitTreeService *CommitTreeService,
	refService *RefService,
	objectService *ObjectService,
) *CommitService {
	return &CommitService{
		writeTreeService:  writeTreeService,
		commitTreeService: commitTreeService,
		refService:        refService,
		objectService:     objectService,
	}
}

func (c *CommitService) Commit(message string) error {
	treeHash, err := c.writeTreeService.WriteTree()
	if err != nil {
		return err
	}

	var parentHashes []string
	headRef, err := c.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}
	parentHash, err := c.refService.Read(headRef)
	if errors.Is(err, ErrRefNotFound) {
		parentHashes = nil
	} else if err != nil {
		return err
	}

	if parentHash != "" {
		parentHashes = append(parentHashes, parentHash)
		parentCommit, err := c.objectService.ReadCommit(parentHash)
		if err != nil {
			return err
		}
		if parentCommit.TreeHash == treeHash {
			return errors.New("nothing to commit")
		}
	}

	commitHash, err := c.commitTreeService.CommitTree(treeHash, message, parentHashes)
	if err != nil {
		return err
	}
	return c.refService.Write(headRef, commitHash)
}
