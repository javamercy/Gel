package core

import "Gel/internal/domain"

type ChangeResult struct {
	IsModified bool
	NewHash    domain.Hash
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

func (c *ChangeDetector) DetectFileChange(entry *domain.IndexEntry, fileStat domain.FileStat) (ChangeResult, error) {
	matches := entry.MatchesStat(fileStat)
	var newHash domain.Hash

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
