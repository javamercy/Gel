package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidAbsolutePath is returned when a path cannot represent an absolute filesystem path.
	ErrInvalidAbsolutePath = errors.New("invalid absolute path")

	// ErrInvalidNormalizedPath is returned when a path is not in repository-normalized form.
	ErrInvalidNormalizedPath = errors.New("invalid normalized path")

	// ErrPathOutsideRepository is returned when a path cannot be represented inside a repository.
	ErrPathOutsideRepository = errors.New("path outside repository")
)

// NormalizedPath represents a repository-relative path using forward slashes.
//
// The root path is represented by the empty string. Non-root paths must not be
// absolute, contain backslashes, contain null bytes, contain empty segments, or
// contain "." / ".." traversal segments.
type NormalizedPath struct {
	value string
}

// NewNormalizedPath resolves path against the current working directory and
// converts it to a repository-relative normalized path.
//
// Use this for command/filesystem input. Use ParseNormalizedPath for paths read
// from .gel data structures.
func NewNormalizedPath(path string, repoDir AbsolutePath) (NormalizedPath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return NormalizedPath{}, fmt.Errorf("normalized path: resolve %q: %w", path, err)
	}
	return absPath.ToNormalizedPath(repoDir)
}

// ParseNormalizedPath validates path as a stored repository-normalized path.
//
// This is intended for paths read from .gel data structures such as the index
// or tree objects. It does not resolve through the current working directory.
func ParseNormalizedPath(path string) (NormalizedPath, error) {
	if err := validateNormalizedFormat(path); err != nil {
		return NormalizedPath{}, err
	}
	return NormalizedPath{value: path}, nil
}

// ToAbsolutePath converts p to an absolute filesystem path under repoDir.
func (n NormalizedPath) ToAbsolutePath(repoDir AbsolutePath) (AbsolutePath, error) {
	if err := validateNormalizedFormat(n.value); err != nil {
		return AbsolutePath{}, err
	}
	if err := repoDir.validate(); err != nil {
		return AbsolutePath{}, fmt.Errorf("normalized path: invalid repository path: %w", err)
	}
	absolutePath := filepath.Join(repoDir.value, filepath.FromSlash(n.value))
	if inside, err := pathWithinDir(repoDir.value, absolutePath); err != nil {
		return AbsolutePath{}, fmt.Errorf("normalized path: compare with repository: %w", err)
	} else if !inside {
		return AbsolutePath{}, fmt.Errorf("%w: %q", ErrPathOutsideRepository, n.value)
	}
	return AbsolutePath{value: absolutePath}, nil
}

// IsRoot reports whether p represents the repository root.
func (n NormalizedPath) IsRoot() bool {
	return n.value == ""
}

// IsWithin reports whether p is equal to or below root.
func (n NormalizedPath) IsWithin(root NormalizedPath) bool {
	if root.IsRoot() {
		return true
	}
	return n.value == root.value || strings.HasPrefix(n.value, root.value+"/")
}

// Equals reports whether p and other are identical.
func (n NormalizedPath) Equals(other NormalizedPath) bool {
	return n == other
}

// String returns the normalized path as a string.
func (n NormalizedPath) String() string {
	return n.value
}

// AbsolutePath represents an absolute filesystem path.
//
// It does not require the path to exist on disk.
type AbsolutePath struct {
	value string
}

// NewAbsolutePath resolves path against the current working directory.
func NewAbsolutePath(path string) (AbsolutePath, error) {
	if path == "" {
		return AbsolutePath{}, fmt.Errorf("%w: empty path", ErrInvalidAbsolutePath)
	}
	if strings.Contains(path, "\x00") {
		return AbsolutePath{}, fmt.Errorf("%w: path contains null byte", ErrInvalidAbsolutePath)
	}

	absPath, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil {
		return AbsolutePath{}, fmt.Errorf("%w: resolve %q: %v", ErrInvalidAbsolutePath, path, err)
	}
	return AbsolutePath{value: absPath}, nil
}

// ToNormalizedPath converts a to a repository-relative normalized path under repoDir.
func (a AbsolutePath) ToNormalizedPath(repoDir AbsolutePath) (NormalizedPath, error) {
	if err := a.validate(); err != nil {
		return NormalizedPath{}, fmt.Errorf("absolute path: normalize %q: %w", a.value, err)
	}
	if err := repoDir.validate(); err != nil {
		return NormalizedPath{}, fmt.Errorf("absolute path: invalid repository path: %w", err)
	}

	relPath, err := filepath.Rel(repoDir.value, a.value)
	if err != nil {
		return NormalizedPath{}, fmt.Errorf("absolute path: make %q relative to %q: %w", a.value, repoDir.value, err)
	}
	if relPath == "." {
		return NormalizedPath{}, nil
	}
	if isRelativeOutside(relPath) {
		return NormalizedPath{}, fmt.Errorf("%w: %q is outside %q", ErrPathOutsideRepository, a.value, repoDir.value)
	}

	normPath := filepath.ToSlash(relPath)

	// TODO: is validation needed here?
	if err := validateNormalizedFormat(normPath); err != nil {
		return NormalizedPath{}, err
	}
	return NormalizedPath{value: normPath}, nil
}

// Equals reports whether p and other are identical.
func (a AbsolutePath) Equals(other AbsolutePath) bool {
	return a == other
}

// String returns the absolute path as a string.
func (a AbsolutePath) String() string {
	return a.value
}

func (a AbsolutePath) validate() error {
	// TODO: we already validated during New, why validate again here? Maybe we can skip validation in ToNormalizedPath since it only accepts AbsolutePaths returned by NewAbsolutePath?
	if a.value == "" {
		return fmt.Errorf("%w: empty path", ErrInvalidAbsolutePath)
	}
	if strings.Contains(a.value, "\x00") {
		return fmt.Errorf("%w: path contains null byte", ErrInvalidAbsolutePath)
	}
	if !filepath.IsAbs(a.value) {
		return fmt.Errorf("%w: %q is not absolute", ErrInvalidAbsolutePath, a.value)
	}
	return nil
}

func validateNormalizedFormat(path string) error {
	if path == "" {
		return nil
	}
	if strings.HasPrefix(path, "/") || filepath.IsAbs(filepath.FromSlash(path)) {
		return fmt.Errorf("%w: %q is absolute", ErrInvalidNormalizedPath, path)
	}
	if strings.Contains(path, "\\") {
		return fmt.Errorf("%w: %q contains backslash", ErrInvalidNormalizedPath, path)
	}
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("%w: path contains null byte", ErrInvalidNormalizedPath)
	}

	segments := strings.Split(path, "/")
	for _, segment := range segments {
		switch segment {
		case "":
			return fmt.Errorf("%w: %q contains empty segment", ErrInvalidNormalizedPath, path)
		case ".", "..":
			return fmt.Errorf("%w: %q contains traversal segment", ErrInvalidNormalizedPath, path)
		}
	}
	return nil
}

func pathWithinDir(repoDir, targetPath string) (bool, error) {
	relPath, err := filepath.Rel(repoDir, targetPath)
	if err != nil {
		return false, err
	}
	return !isRelativeOutside(relPath), nil
}

func isRelativeOutside(path string) bool {
	return path == ".." || strings.HasPrefix(path, ".."+string(filepath.Separator))
}
