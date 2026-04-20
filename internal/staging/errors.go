package staging

import (
	"Gel/internal/domain"
	"errors"
	"fmt"
)

var (
	// ErrPathDidNotMatch is returned when a pathspec matches no files in the index or working tree.
	ErrPathDidNotMatch = errors.New("pathspec did not match any files")

	errRemovePathDidNotMatch        = errors.New("remove: pathspec did not match any tracked files")
	errRemoveRecursiveRequired      = errors.New("remove: recursive flag required")
	errRemoveOutsideRepository      = errors.New("remove: path is outside repository")
	errRemoveHasStagedChanges       = errors.New("remove: file has staged changes")
	errRemoveHasLocalModifications  = errors.New("remove: file has local modifications")
	errRemoveHasStagedAndLocalState = errors.New("remove: file has staged and local changes")
)

type removePathDidNotMatchError struct {
	pathspec string
}

func (e *removePathDidNotMatchError) Error() string {
	return fmt.Sprintf("fatal: pathspec '%s' did not match any files", e.pathspec)
}

func (e *removePathDidNotMatchError) Unwrap() error {
	return errRemovePathDidNotMatch
}

type removeRecursiveRequiredError struct {
	pathspec string
}

func (e *removeRecursiveRequiredError) Error() string {
	return fmt.Sprintf("fatal: not removing '%s' recursively without -r", e.pathspec)
}

func (e *removeRecursiveRequiredError) Unwrap() error {
	return errRemoveRecursiveRequired
}

type removeOutsideRepositoryError struct {
	pathspec string
}

func (e *removeOutsideRepositoryError) Error() string {
	return fmt.Sprintf("fatal: '%s' is outside repository", e.pathspec)
}

func (e *removeOutsideRepositoryError) Unwrap() error {
	return errRemoveOutsideRepository
}

type removeHasStagedChangesError struct {
	path domain.NormalizedPath
}

func (e *removeHasStagedChangesError) Error() string {
	return fmt.Sprintf(
		"error: the following file has changes staged in the index:\n    %s\n(use --cached to keep the file, or -f to force removal)",
		e.path,
	)
}

func (e *removeHasStagedChangesError) Unwrap() error {
	return errRemoveHasStagedChanges
}

type removeHasLocalModificationsError struct {
	path domain.NormalizedPath
}

func (e *removeHasLocalModificationsError) Error() string {
	return fmt.Sprintf(
		"error: the following file has local modifications:\n    %s\n(use --cached to keep the file, or -f to force removal)",
		e.path,
	)
}

func (e *removeHasLocalModificationsError) Unwrap() error {
	return errRemoveHasLocalModifications
}

type removeHasStagedAndLocalStateError struct {
	path domain.NormalizedPath
}

func (e *removeHasStagedAndLocalStateError) Error() string {
	return fmt.Sprintf(
		"error: the following file has staged content different from both the file and the HEAD:\n    %s\n(use -f to force removal)",
		e.path,
	)
}

func (e *removeHasStagedAndLocalStateError) Unwrap() error {
	return errRemoveHasStagedAndLocalState
}

func newRemovePathDidNotMatchError(pathspec string) error {
	return &removePathDidNotMatchError{pathspec: pathspec}
}

func newRemoveRecursiveRequiredError(pathspec string) error {
	return &removeRecursiveRequiredError{pathspec: pathspec}
}

func newRemoveOutsideRepositoryError(pathspec string) error {
	return &removeOutsideRepositoryError{pathspec: pathspec}
}

func newRemoveHasStagedChangesError(path domain.NormalizedPath) error {
	return &removeHasStagedChangesError{path: path}
}

func newRemoveHasLocalModificationsError(path domain.NormalizedPath) error {
	return &removeHasLocalModificationsError{path: path}
}

func newRemoveHasStagedAndLocalStateError(path domain.NormalizedPath) error {
	return &removeHasStagedAndLocalStateError{path: path}
}
