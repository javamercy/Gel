package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RootPath represents an empty normalized path, used for the repository root.
const RootPath = NormalizedPath("")

// NormalizedPath represents a path relative to the repository root,
// using forward slashes (e.g., "src/main.go").
// It cannot be absolute, contain backslashes, or contain null bytes.
type NormalizedPath string

// NewNormalizedPath creates a NormalizedPath from an arbitrary path string
// relative to the repository directory.
func NewNormalizedPath(repoDir string, path string) (NormalizedPath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path: %w", err)
	}
	return absPath.ToNormalizedPath(repoDir)
}

// NewNormalizedPathUnchecked creates a NormalizedPath without converting through
// the filesystem. It validates that the path is in normalized format.
func NewNormalizedPathUnchecked(path string) (NormalizedPath, error) {
	if err := validateNormalizedFormat(path); err != nil {
		return "", err
	}
	return NormalizedPath(path), nil
}

// ToAbsolutePath converts a normalized path to an absolute path within the repository.
func (p NormalizedPath) ToAbsolutePath(repoDir string) (AbsolutePath, error) {
	absPath := filepath.Join(repoDir, filepath.FromSlash(p.String()))
	return AbsolutePath(absPath), nil
}

// Equals reports whether two normalized paths are identical.
func (p NormalizedPath) Equals(o NormalizedPath) bool {
	return p == o
}

func (p NormalizedPath) String() string {
	return string(p)
}

// validateNormalizedFormat checks that a path is in valid normalized format
// (no leading slash, no backslashes, no null bytes).
func validateNormalizedFormat(path string) error {
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("normalized path cannot be absolute: %s", path)
	}
	if strings.Contains(path, "\\") {
		return fmt.Errorf("normalized path must use forward slashes: %s", path)
	}
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("normalized path cannot contain null bytes: %s", path)
	}
	return nil
}

// AbsolutePath represents an absolute filesystem path (e.g., "/home/user/project/src/main.go").
type AbsolutePath string

// NewAbsolutePath creates an AbsolutePath from an arbitrary path string.
// It normalizes the path to use forward slashes internally.
func NewAbsolutePath(path string) (AbsolutePath, error) {
	absPath, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return AbsolutePath(absPath), nil
}

// ToNormalizedPath converts an absolute path to a normalized path relative
// to the repository directory. Returns RootPath if the absolute path is exactly
// the repository root.
func (p AbsolutePath) ToNormalizedPath(repoDir string) (NormalizedPath, error) {
	relPath, err := filepath.Rel(repoDir, p.String())
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	normPath := filepath.ToSlash(relPath)
	if normPath == "." {
		return RootPath, nil
	}
	if err := validateNormalizedFormat(normPath); err != nil {
		return "", fmt.Errorf("path is outside repository or invalid: %w", err)
	}
	return NormalizedPath(normPath), nil
}

func (p AbsolutePath) String() string {
	return string(p)
}
