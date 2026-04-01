package core

import (
	domain2 "Gel/internal/domain"
)

type ChangeResult struct {
	IsModified bool
	NewHash    domain2.Hash
}

type ChangeDetector struct {
	objectService *ObjectService
	repoDir       string
}

func NewChangeDetector(objectService *ObjectService, repoDir string) *ChangeDetector {
	return &ChangeDetector{
		objectService: objectService,
		repoDir:       repoDir,
	}
}

func (c *ChangeDetector) DetectFileChange(entry *domain2.IndexEntry, fileStat domain2.FileStat) (ChangeResult, error) {
	matches := entry.MatchesStat(fileStat)
	var newHash domain2.Hash

	if !matches {
		absolutePath, err := entry.Path.ToAbsolutePath(c.repoDir)
		if err != nil {
			return ChangeResult{}, err
		}
		newHash, _, err = c.objectService.ComputeObjectHash(absolutePath)
		if err != nil {
			return ChangeResult{}, err
		}
	}
	return ChangeResult{
		IsModified: !matches,
		NewHash:    newHash,
	}, nil
}
