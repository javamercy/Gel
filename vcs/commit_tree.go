package vcs

import (
	"Gel/core/encoding"
	"Gel/core/util"
	"Gel/core/validation"
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

func (commitTreeService *CommitTreeService) CommitTree(hash string, message string) (string, error) {

	if err := validation.ValidateHash(hash); err != nil {
		return "", err
	}

	object, err := commitTreeService.objectService.Read(hash)
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
		TreeHash:     hash,
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
	commitHash := encoding.ComputeSha256(serializedCommit)
	err = commitTreeService.objectService.Write(commitHash, serializedCommit)
	if err != nil {
		return "", err
	}
	return hash, nil
}
