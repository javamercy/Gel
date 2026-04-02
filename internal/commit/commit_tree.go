package commit

import (
	"Gel/internal/core"
	"Gel/internal/domain"
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

func (c *CommitTreeService) CommitTree(hash domain.Hash, message string, parentHashes []domain.Hash) (
	domain.Hash, error,
) {
	_, err := c.objectService.ReadTree(hash)
	if err != nil {
		return domain.Hash{}, err
	}

	name, email, err := c.configService.GetUserInfo()
	if err != nil {
		return domain.Hash{}, err
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
	hexCommitHash := core.ComputeSHA256(serializedData)
	commitHash, err := domain.NewHash(hexCommitHash)
	if err != nil {
		return domain.Hash{}, err
	}
	if err := c.objectService.Write(commitHash, serializedData); err != nil {
		return domain.Hash{}, err
	}
	return commitHash, nil
}
