package vcs

import (
	"Gel/core/encoding"
	"Gel/domain"
)

type CommitTreeService struct {
	objectService *ObjectService
}

func NewCommitTreeService(objectService *ObjectService) *CommitTreeService {
	return &CommitTreeService{
		objectService: objectService,
	}
}

func (commitTreeService *CommitTreeService) CommitTree(treeHash string, message string) (string, error) {

	object, err := commitTreeService.objectService.Read(treeHash)
	if err != nil {
		return "", err
	}
	_, ok := object.(*domain.Tree)
	if !ok {
		return "", domain.ErrInvalidObjectType
	}

	commitFields := domain.CommitFields{
		TreeHash:     treeHash,
		ParentHashes: nil,
		Author: domain.Identity{
			Name:      "Linus Torvalds",
			Email:     "torvalds@linux-foundation.org",
			Timestamp: "2025-01-01T00:00:00Z",
			Timezone:  "+0000",
		},
		Committer: domain.Identity{
			Name:      "Linus Torvalds",
			Email:     "torvalds@linux-foundation.org",
			Timestamp: "2025-01-01T00:00:00Z",
			Timezone:  "+0000",
		},
		Message: message,
	}

	commit := domain.NewCommitFromFields(commitFields)
	serializedCommit := commit.Serialize()
	hash := encoding.ComputeHash(serializedCommit)
	err = commitTreeService.objectService.Write(hash, serializedCommit)
	if err != nil {
		return "", err
	}

	return hash, nil
}
