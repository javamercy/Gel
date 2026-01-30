package gel

import (
	"Gel/internal/workspace"
	"errors"
	"io/fs"
)

type CommitService struct {
	writeTreeService  *WriteTreeService
	commitTreeService *CommitTreeService
	refService        *RefService
	objectService     *ObjectService
}

func NewCommitService(writeTreeService *WriteTreeService,
	commitTreeService *CommitTreeService,
	refService *RefService,
	objectService *ObjectService) *CommitService {
	return &CommitService{
		writeTreeService:  writeTreeService,
		commitTreeService: commitTreeService,
		refService:        refService,
		objectService:     objectService,
	}
}

func (s *CommitService) Commit(message string) error {
	treeHash, err := s.writeTreeService.WriteTree()
	if err != nil {
		return err
	}

	var parentHashes []string
	headRef, err := s.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}
	parentHash, err := s.refService.Read(headRef)
	if errors.Is(err, fs.ErrNotExist) {
		parentHashes = nil
	} else if err != nil {
		return err
	}

	if parentHash != "" {
		parentHashes = append(parentHashes, parentHash)
		parentCommit, err := s.objectService.ReadCommit(parentHash)
		if err != nil {
			return err
		}
		if parentCommit.TreeHash == treeHash {
			return errors.New("nothing to commit")
		}
	}

	commitHash, err := s.commitTreeService.CommitTree(treeHash, message, parentHashes)
	if err != nil {
		return err
	}
	return s.refService.Write(headRef, commitHash)
}
