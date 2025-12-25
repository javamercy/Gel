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

	author, err := domain.NewIdentity(
		"Linus Torvalds",
		"torvalds@linux-foundation.org",
		"2025-01-01T00:00:00Z",
		"+0000")
	if err != nil {
		return "", err
	}

	commitFields := domain.CommitFields{
		TreeHash:     treeHash,
		ParentHashes: nil,
		Author:       author,
		Committer:    author,
		Message:      message,
	}

	commit, err := domain.NewCommitFromFields(commitFields)
	if err != nil {
		return "", err
	}
	serializedCommit := commit.Serialize()
	hash := encoding.ComputeHash(serializedCommit)
	err = commitTreeService.objectService.Write(hash, serializedCommit)
	if err != nil {
		return "", err
	}

	return hash, nil
}
