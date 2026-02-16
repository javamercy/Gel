package commit

import (
	"Gel/domain"
	core2 "Gel/internal/core"
	"time"
)

type CommitTreeService struct {
	objectService *core2.ObjectService
	configService *core2.ConfigService
}

func NewCommitTreeService(objectService *core2.ObjectService, configService *core2.ConfigService) *CommitTreeService {
	return &CommitTreeService{
		objectService: objectService,
		configService: configService,
	}
}

func (c *CommitTreeService) CommitTree(hash string, message string, parentHashes []string) (string, error) {
	_, err := c.objectService.ReadTree(hash)
	if err != nil {
		return "", err
	}

	name, email, err := c.configService.GetUserInfo()
	if err != nil {
		return "", err
	}

	now := time.Now()

	identity := domain.NewIdentity(
		name,
		email,
		domain.FormatCommitTimestamp(now),
		domain.FormatCommitTimezone(now),
	)
	commitFields := domain.CommitFields{
		TreeHash:     hash,
		ParentHashes: parentHashes,
		Author:       identity,
		Committer:    identity,
		Message:      message,
	}
	commit := domain.NewCommitFromFields(commitFields)
	serializedData := commit.Serialize()
	commitHash := core2.ComputeSHA256(serializedData)
	err = c.objectService.Write(commitHash, serializedData)
	if err != nil {
		return "", err
	}
	return commitHash, nil
}
