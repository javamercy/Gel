package vcs

import (
	"Gel/core/encoding"
	"Gel/core/util"
	"Gel/domain"
	"time"
)

type CommitTreeService struct {
	objectService *ObjectService
	configService *ConfigService
}

func NewCommitTreeService(objectService *ObjectService, configService *ConfigService) *CommitTreeService {
	return &CommitTreeService{
		objectService: objectService,
		configService: configService,
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

	userIdentity, err := commitTreeService.configService.GetUserIdentity()
	if err != nil {
		return "", err
	}
	now := time.Now()

	author, err := domain.NewIdentity(
		userIdentity.Name,
		userIdentity.Email,
		util.FormatCommitTimestamp(now),
		util.FormatCommitTimezone(now),
	)

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
