package core

import "errors"

var (
	// ErrRefNotFound is returned when a reference (branch, tag, HEAD) does not exist on disk.
	ErrRefNotFound = errors.New("reference not found")

	// ErrInvalidRef is returned when a ref path does not start with "refs/".
	ErrInvalidRef = errors.New("invalid ref: must start with refs/")

	// ErrInvalidSymbolicRef is returned when a symbolic ref file has a malformed format.
	ErrInvalidSymbolicRef = errors.New("invalid symbolic ref")

	// ErrRefUpdateConflict is returned when a safe update finds the ref points to an unexpected hash.
	ErrRefUpdateConflict = errors.New("ref update conflict: current hash does not match expected")

	// ErrPathNotFoundInTree is returned when a path lookup in a tree object finds no match.
	ErrPathNotFoundInTree = errors.New("path not found in tree")

	// ErrUnknownPathspecType is returned when a pathspec cannot be classified.
	ErrUnknownPathspecType = errors.New("unknown pathspec type")
)
