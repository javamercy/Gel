package commit

import (
	"Gel/domain"
	"Gel/internal/core"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

type LogService struct {
	refService    *core.RefService
	objectService *core.ObjectService
}

func NewLogService(refService *core.RefService, objectService *core.ObjectService) *LogService {
	return &LogService{
		refService:    refService,
		objectService: objectService,
	}
}

func (l *LogService) Log(writer io.Writer, name string, limit int, oneline bool) error {
	hash, err := l.refService.Resolve(name)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("current branch '%s' does not have any commits", name)
	}

	count := 0
	for hash != "" {
		if limit > 0 && count >= limit {
			break
		}
		count++
		commit, err := l.objectService.ReadCommit(hash)
		if err != nil {
			return err
		}
		if oneline {
			if err := l.printCommitOneline(writer, hash, commit); err != nil {
				return err
			}
		} else {
			if err := l.printCommit(writer, hash, commit); err != nil {
				return err
			}
		}

		if len(commit.ParentHashes) == 0 {
			break
		}
		// TODO: use a priority queue/heap to interleave commits from all parents by date
		hash = commit.ParentHashes[0]
	}
	return nil
}

func (l *LogService) printCommit(writer io.Writer, hash string, commit *domain.Commit) error {
	t, err := domain.FormatCommitDate(commit.Author.Timestamp, commit.Author.Timezone)
	if err != nil {
		return err
	}
	commitPrefix := core.ColorGreen
	commitSuffix := core.ColorReset
	if _, err := fmt.Fprintf(
		writer,
		"%scommit %s%s\n"+
			"Author: %s <%s>\n"+
			"Date:   %v\n"+
			"\n    %s\n\n",
		commitPrefix,
		hash,
		commitSuffix,
		commit.Author.Name, commit.Author.Email,
		t,
		commit.Message,
	); err != nil {
		return err
	}
	return nil
}

func (l *LogService) printCommitOneline(writer io.Writer, hash string, commit *domain.Commit) error {
	shortHash := hash[:7]
	commitPrefix := core.ColorGreen
	commitSuffix := core.ColorReset
	firstLine := strings.Split(commit.Message, "\n")[0]
	if _, err := fmt.Fprintf(writer, "%s%s%s %s\n", commitPrefix, shortHash, commitSuffix, firstLine); err != nil {
		return err
	}
	return nil
}
