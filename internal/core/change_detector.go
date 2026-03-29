package core

import (
	"Gel/domain"
)

type ChangeResult struct {
	IsModified bool
	NewHash    domain.Hash
}

type ChangeDetector struct {
	objectService *ObjectService
}

func NewChangeDetector(objectService *ObjectService) *ChangeDetector {
	return &ChangeDetector{
		objectService: objectService,
	}
}

func (c *ChangeDetector) DetectFileChange(entry *domain.IndexEntry, fileStat domain.FileStat) (ChangeResult, error) {
	matches := entry.MatchesStat(fileStat)
	var newHash domain.Hash
	var err error

	if !matches {
		newHash, _, err = c.objectService.ComputeObjectHash(entry.Path.ToAbsolutePath())
	}
	if err != nil {
		return ChangeResult{}, err
	}
	return ChangeResult{
		IsModified: !matches,
		NewHash:    newHash,
	}, nil
}
