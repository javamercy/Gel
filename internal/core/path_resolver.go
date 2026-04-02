package core

import (
	"Gel/internal/domain"
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
		// TODO: implement gelignore.
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
		pathspecType, err := classifyPathspec(pathspec)
		if err != nil {
			return nil, err
		}

		paths, err := p.expandPathspec(pathspec, pathspecType)
		if err != nil {
			return nil, err
		}

		normalizedScope, err := p.normalizeScope(pathspec, pathspecType)
		if err != nil {
			return nil, err
		}

		normalizedPaths := make(map[domain.NormalizedPath]bool)
		for _, path := range paths {
			normalizedPath, err := domain.NewNormalizedPath(p.repoDir, path)
			if err != nil {
				return nil, err
			}
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
	segments := strings.Split(filepath.ToSlash(path), "/")
	for _, segment := range segments {
		if p.ignoredPatterns[segment] {
			return true
		}
	}
	return false
}

func (p *PathResolver) expandPathspec(pathspec string, pathspecType PathspecType) ([]string, error) {
	switch pathspecType {
	case PathspecTypeFile:
		return []string{pathspec}, nil
	case PathspecTypeDirectory:
		return p.expandDirectory(pathspec)
	case PathspecTypeGlobPattern:
		return p.expandGlobPattern(pathspec)
	case PathspecTypeNonExistent:
		return []string{}, nil
	default:
		return nil, ErrUnknownPathspecType
	}
}

func (p *PathResolver) normalizeScope(pathspec string, pathspecType PathspecType) (string, error) {
	switch pathspecType {
	case PathspecTypeFile,
		PathspecTypeDirectory,
		PathspecTypeGlobPattern,
		PathspecTypeNonExistent:
		np, err := domain.NewNormalizedPath(p.repoDir, pathspec)
		if err != nil {
			return "", err
		}
		return np.String(), nil
	default:
		return "", ErrUnknownPathspecType
	}
	if normalizedScope == "." {
		normalizedScope = ""
	}
	return normalizedScope, nil
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

func (p *PathResolver) expandDirectory(path string) ([]string, error) {
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

func (p *PathResolver) expandGlobPattern(pattern string) ([]string, error) {
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
			files, err := p.expandDirectory(match)
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
