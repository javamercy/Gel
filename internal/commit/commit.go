package commit

import (
	core2 "Gel/internal/core"
	"Gel/internal/tree"
	"Gel/internal/workspace"
	"errors"
)

type CommitService struct {
	writeTreeService  *tree.WriteTreeService
	commitTreeService *CommitTreeService
	refService        *core2.RefService
	objectService     *core2.ObjectService
}

func NewCommitService(
	writeTreeService *tree.WriteTreeService,
	commitTreeService *CommitTreeService,
	refService *core2.RefService,
	objectService *core2.ObjectService,
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
	if errors.Is(err, core2.ErrRefNotFound) {
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
