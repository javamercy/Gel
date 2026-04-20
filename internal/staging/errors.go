package staging

import (
	"Gel/internal/domain"
	"errors"
	"fmt"
)

var (
	// ErrPathDidNotMatch is returned when a pathspec matches no files in the index or working tree.
	ErrPathDidNotMatch = errors.New("pathspec did not match any files")

	// errRemovePathDidNotMatch identifies rm failures where no tracked path matches the pathspec.
	errRemovePathDidNotMatch = errors.New("remove: pathspec did not match any tracked files")

	// errRemoveRecursiveRequired identifies rm failures where a directory-like pathspec is used without -r.
	errRemoveRecursiveRequired = errors.New("remove: recursive flag required")

	// errRemoveOutsideRepository identifies rm failures where the requested path resolves outside the repo root.
	errRemoveOutsideRepository = errors.New("remove: path is outside repository")

	// errRemoveHasStagedChanges identifies safety failures where staged changes would be discarded.
	errRemoveHasStagedChanges = errors.New("remove: file has staged changes")

	// errRemoveHasLocalModifications identifies safety failures where working tree changes would be discarded.
	errRemoveHasLocalModifications = errors.New("remove: file has local modifications")

	// errRemoveHasStagedAndLocalState identifies safety failures where both staged and local changes exist.
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
