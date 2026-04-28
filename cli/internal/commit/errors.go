package commit

import "errors"

var (
	// ErrNothingToCommit is returned when the tree matches the parent commit.
	ErrNothingToCommit = errors.New("nothing to commit")

	// ErrNoCommitsYet is returned when trying to log a branch with no commits.
	ErrNoCommitsYet = errors.New("no commits yet")
)
