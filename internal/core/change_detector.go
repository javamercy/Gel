package core

import "Gel/domain"

type ChangeResult struct {
	IsModified bool
	NewHash    string
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
	var newHash string
	var err error

	if !matches {
		newHash, _, err = c.hashObjectService.HashObject(entry.Path, false)
	}
	if err != nil {
		return ChangeResult{}, err
	}
	return ChangeResult{
		IsModified: !matches,
		NewHash:    newHash,
	}, nil
}
