package commit

import (
	"Gel/domain"
	"Gel/internal/gel/core"
	"time"
)

type CommitTreeService struct {
	objectService *core.ObjectService
	configService *core.ConfigService
}

func NewCommitTreeService(objectService *core.ObjectService, configService *core.ConfigService) *CommitTreeService {
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
	commitHash := core.ComputeSHA256(serializedData)
	err = c.objectService.Write(commitHash, serializedData)
	if err != nil {
		return "", err
	}
	return commitHash, nil
}
