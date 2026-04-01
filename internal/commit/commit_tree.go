package commit

import (
	"Gel/internal/core"
	domain2 "Gel/internal/domain"
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

func (c *CommitTreeService) CommitTree(hash domain2.Hash, message string, parentHashes []domain2.Hash) (
	domain2.Hash, error,
) {
	_, err := c.objectService.ReadTree(hash)
	if err != nil {
		return domain2.Hash{}, err
	}

	name, email, err := c.configService.GetUserInfo()
	if err != nil {
		return domain2.Hash{}, err
	}

	now := time.Now()

	identity := domain2.NewIdentity(
		name,
		email,
		domain2.FormatCommitTimestamp(now),
		domain2.FormatCommitTimezone(now),
	)
	commitFields := domain2.CommitFields{
		TreeHash:     hash,
		ParentHashes: parentHashes,
		Author:       identity,
		Committer:    identity,
		Message:      message,
	}
	commit := domain2.NewCommitFromFields(commitFields)
	serializedData := commit.Serialize()
	hexCommitHash := core.ComputeSHA256(serializedData)
	commitHash, err := domain2.NewHash(hexCommitHash)
	if err != nil {
		return domain2.Hash{}, err
	}
	if err := c.objectService.Write(commitHash, serializedData); err != nil {
		return domain2.Hash{}, err
	}
	return commitHash, nil
}
