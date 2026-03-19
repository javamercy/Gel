package core

import "Gel/domain"

type ChangeResult struct {
	IsModified bool
	NewHash    domain.Hash
}

type ChangeDetector struct {
	hashObjectService *HashObjectService
}

func NewChangeDetector(hashObjectService *HashObjectService) *ChangeDetector {
	return &ChangeDetector{
		hashObjectService: hashObjectService,
	}
}

func (c *ChangeDetector) DetectFileChange(
	entry *domain.IndexEntry, fileStat domain.FileStat,
) (ChangeResult, error) {
	matches := entry.MatchesStat(fileStat)
	var newHash domain.Hash
	var err error

	if !matches {
		newHash, _, err = c.hashObjectService.ComputeObjectHash(entry.Path)
	}
	if err != nil {
		return ChangeResult{}, err
	}
	return ChangeResult{
		IsModified: !matches,
		NewHash:    newHash,
	}, nil
}
