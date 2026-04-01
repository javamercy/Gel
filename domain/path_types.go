package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// NormalizedPath represents a path relative to repository root with forward slashes.
type NormalizedPath string

func NewNormalizedPath(repoDir string, path string) (NormalizedPath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return "", fmt.Errorf("failed to create absolute path: %w", err)
	}
	return absPath.ToNormalizedPath(repoDir)
}

// NewNormalizedPathUnchecked creates a NormalizedPath from an already-normalized string
// (e.g., from index file or other storage). It validates the format but doesn't convert.
// Use this when you know the path is already in normalized format.
func NewNormalizedPathUnchecked(path string) (NormalizedPath, error) {
	if err := validateNormalizedFormat(path); err != nil {
		return "", err
	}
	return NormalizedPath(path), nil
}

func (p NormalizedPath) ToAbsolutePath(repoDir string) (AbsolutePath, error) {
	absPath := filepath.Join(repoDir, filepath.FromSlash(p.String()))
	return AbsolutePath(absPath), nil
}

func (p NormalizedPath) String() string {
	return string(p)
}

// validateNormalizedFormat ensures a path is in proper normalized format.
// This provides defense against corrupted index files or malicious data.
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

// AbsolutePath represents an absolute path to a file or directory.
type AbsolutePath string

func NewAbsolutePath(path string) (AbsolutePath, error) {
	absPath, err := filepath.Abs(filepath.FromSlash(path))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return AbsolutePath(absPath), nil
}

func (p AbsolutePath) ToNormalizedPath(repoDir string) (NormalizedPath, error) {
	relPath, err := filepath.Rel(repoDir, p.String())
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	normalizedPath := filepath.ToSlash(relPath)
	if normalizedPath == "." {
		normalizedPath = ""
	}
	if err := validateNormalizedFormat(normalizedPath); err != nil {
		return "", fmt.Errorf("path is outside repository or invalid: %w", err)
	}
	return NormalizedPath(normalizedPath), nil
}

func (p AbsolutePath) String() string {
	return string(p)
}
