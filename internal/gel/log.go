package gel

import (
	"Gel/domain"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

type LogService struct {
	refService    *RefService
	objectService *ObjectService
}

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func NewLogService(refService *RefService, objectService *ObjectService) *LogService {
	return &LogService{
		refService:    refService,
		objectService: objectService,
	}
}

func (s *LogService) Log(w io.Writer, name string, limit int, oneline bool) error {
	hash, err := s.refService.Resolve(name)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("current branch '%s' does not have any commits", name)
	}

	count := 0
	for hash != "" {
		if limit > 0 && count >= limit {
			break
		}
		count++
		commit, err := s.objectService.ReadCommit(hash)
		if err != nil {
			return err
		}
		if oneline {
			if err := s.printCommitOneline(w, hash, commit); err != nil {
				return err
			}
		} else {
			if err := s.printCommit(w, hash, commit); err != nil {
				return err
			}
		}

		if len(commit.ParentHashes) == 0 {
			break
		}
		hash = commit.ParentHashes[0]
	}
	return nil
}

func (s *LogService) printCommit(w io.Writer, hash string, commit *domain.Commit) error {
	t, err := domain.FormatCommitDate(commit.Author.Timestamp, commit.Author.Timezone)
	if err != nil {
		return err
	}
	commitPrefix := colorGreen
	commitSuffix := colorReset
	if _, err := fmt.Fprintf(w,
		"%scommit %s%s\n"+
			"Author: %s <%s>\n"+
			"Date:   %v\n"+
			"\n    %s\n\n",
		commitPrefix,
		hash,
		commitSuffix,
		commit.Author.User.Name, commit.Author.User.Email,
		t,
		commit.Message); err != nil {
		return err
	}
	return nil
}

func (s *LogService) printCommitOneline(w io.Writer, hash string, commit *domain.Commit) error {
	shortHash := hash[:7]
	commitPrefix := colorGreen
	commitSuffix := colorReset
	firstLine := strings.Split(commit.Message, "\n")[0]
	if _, err := fmt.Fprintf(w, "%s%s%s %s\n", commitPrefix, shortHash, commitSuffix, firstLine); err != nil {
		return err
	}
	return nil
}
