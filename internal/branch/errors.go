package branch

import "errors"

var (
	// ErrBranchNotFound is returned when a branch does not exist.
	ErrBranchNotFound = errors.New("branch not found")

	// ErrBranchAlreadyExists is returned when trying to create a branch that already exists.
	ErrBranchAlreadyExists = errors.New("branch already exists")

	// ErrDeleteCurrentBranch is returned when trying to delete the currently checked-out branch.
	ErrDeleteCurrentBranch = errors.New("cannot delete the current branch")

	// ErrInvalidBranchName is returned when a branch name violates naming rules.
	ErrInvalidBranchName = errors.New("invalid branch name")

	// ErrUncommittedChanges is returned when switching branches with staged local changes.
	ErrUncommittedChanges = errors.New("uncommitted changes")
)
