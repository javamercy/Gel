package gel

import (
	"Gel/domain"
	"Gel/internal/gel/validate"
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

func (commitTreeService *CommitTreeService) CommitTree(hash string, message string, parentHashes []string) (string, error) {
	if err := validate.Hash(hash); err != nil {
		return "", err
	}

	_, err := commitTreeService.objectService.ReadTree(hash)
	if err != nil {
		return "", err
	}

	user, err := commitTreeService.configService.GetUserIdentity()
	if err != nil {
		return "", err
	}

	now := time.Now()

	identity, err := domain.NewIdentity(
		user.Name,
		user.Email,
		domain.FormatCommitTimestamp(now),
		domain.FormatCommitTimezone(now),
	)
	if err != nil {
		return "", err
	}

	commitFields := domain.CommitFields{
		TreeHash:     hash,
		ParentHashes: parentHashes,
		Author:       identity,
		Committer:    identity,
		Message:      message,
	}
	commit, err := domain.NewCommitFromFields(commitFields)
	if err != nil {
		return "", err
	}

	serializedData := commit.Serialize()
	commitHash := ComputeSHA256(serializedData)
	err = commitTreeService.objectService.Write(commitHash, serializedData)
	if err != nil {
		return "", err
	}
	return commitHash, nil
}
