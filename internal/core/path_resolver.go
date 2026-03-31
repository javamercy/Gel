package core

import (
	"Gel/domain"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	PathspecTypeFile PathspecType = iota
	PathspecTypeDirectory
	PathspecTypeGlobPattern
	PathspecTypeNonExistent
)

const (
	globPatterns string = "*?[]"
)

type PathspecType int

func (t PathspecType) String() string {
	switch t {
	case PathspecTypeFile:
		return "file"
	case PathspecTypeDirectory:
		return "directory"
	case PathspecTypeGlobPattern:
		return "glob pattern"
	case PathspecTypeNonExistent:
		return "non-existent"
	}
	return "unknown"
}

type ResolvedPath struct {
	Type            PathspecType
	NormalizedScope string
	NormalizedPaths map[domain.NormalizedPath]bool
}

type PathResolver struct {
	repoDir         string
	ignoredPatterns map[string]bool
}

func NewPathResolver(repoDir string, ignoredPatterns map[string]bool) *PathResolver {
	if ignoredPatterns == nil {
		ignoredPatterns = map[string]bool{
			".gel":  true,
			".git":  true,
			".idea": true,
		}
	}
	return &PathResolver{
		repoDir:         repoDir,
		ignoredPatterns: ignoredPatterns,
	}
}

func (p *PathResolver) Resolve(pathspecs []string) ([]ResolvedPath, error) {
	resolvedPaths := make([]ResolvedPath, 0)

	for _, pathspec := range pathspecs {
		var paths []string
		var err error

		pathspecType, err := classifyPathspec(pathspec)
		if err != nil {
			return nil, err
		}
		absPathspec, err := domain.NewAbsolutePath(pathspec)
		if err != nil {
			return nil, err
		}

		normalizedPath, err := absPathspec.ToNormalizedPath(p.repoDir)
		if err != nil {
			return nil, err
		}

		normalizedScope := normalizedPath.String()
		if normalizedScope == "." {
			normalizedScope = ""
		}
		switch pathspecType {
		case PathspecTypeFile:
			paths = []string{pathspec}
		case PathspecTypeDirectory:
			paths, err = expandDirectory(pathspec)
		case PathspecTypeGlobPattern:
			paths, err = expandGlobPattern(pathspec)
		case PathspecTypeNonExistent:
			paths = []string{}
		default:
			return nil, ErrUnknownPathspecType
		}

		if err != nil {
			return nil, err
		}

		normalizedPaths := make(map[domain.NormalizedPath]bool)
		for _, path := range paths {
			normalizedPath, err := domain.NewNormalizedPath(p.repoDir, path)
			if err != nil {
				return nil, err
			}
			// TODO: implement gelignore.
			if p.shouldIgnore(normalizedPath.String()) || normalizedPaths[normalizedPath] {
				continue
			}
			normalizedPaths[normalizedPath] = true
		}

		resolvedPaths = append(
			resolvedPaths, ResolvedPath{
				Type:            pathspecType,
				NormalizedScope: normalizedScope,
				NormalizedPaths: normalizedPaths,
			},
		)
	}
	return resolvedPaths, nil
}

func (p *PathResolver) shouldIgnore(path string) bool {
	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if p.ignoredPatterns[segment] {
			return true
		}
	}
	return false
}

func classifyPathspec(pathspec string) (PathspecType, error) {
	if strings.ContainsAny(pathspec, globPatterns) {
		return PathspecTypeGlobPattern, nil
	}

	fileInfo, err := os.Stat(pathspec)
	if errors.Is(err, os.ErrNotExist) {
		return PathspecTypeNonExistent, nil
	} else if err != nil {
		return PathspecTypeNonExistent, err
	}
	if fileInfo.IsDir() {
		return PathspecTypeDirectory, nil
	}
	return PathspecTypeFile, nil
}

func expandDirectory(path string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(
		path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				info, err := os.Stat(p)
				if err == nil && info.IsDir() {
					return nil
				}
				files = append(files, p)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func expandGlobPattern(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.IsDir() {
			files, err := expandDirectory(match)
			if err != nil {
				return nil, err
			}
			result = append(result, files...)
		} else {
			result = append(result, match)
		}
	}
	return result, nil
}
