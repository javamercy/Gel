package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// NormalizedPath represents a path relative to repository root with forward slashes.
type NormalizedPath string

func NewNormalizedPath(path string) (NormalizedPath, error) {
	if err := validateNormalizedFormat(path); err != nil {
		return "", err
	}
	return NormalizedPath(path), nil
}

func NewNormalizedPathFromAbsolutePath(path string) (NormalizedPath, error) {
	absPath, err := NewAbsolutePath(path)
	if err != nil {
		return "", err
	}
	return absPath.ToNormalizedPath(filepath.Dir(path))
}

func validateNormalizedFormat(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("normalized path cannot be absolute: %s", path)
	}
	if strings.Contains(path, "\\") {
		return fmt.Errorf("normalized path must use forward slashes: %s", path)
	}

	components := strings.Split(path, "/")
	for _, comp := range components {
		if comp == "." || comp == ".." {
			return fmt.Errorf("normalized path cannot contain %s: %s", comp, path)
		}
	}
	return nil
}

func (p NormalizedPath) ToAbsolutePath() (AbsolutePath, error) {
	return NewAbsolutePath(p.String())
}

func (p NormalizedPath) String() string {
	return string(p)
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
	return NormalizedPath(filepath.ToSlash(relPath)), nil
}

func (p AbsolutePath) String() string {
	return string(p)
}
