package staging

import (
	"Gel/internal/core"
	"Gel/internal/domain"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RemoveOptions controls rm command execution mode.
type RemoveOptions struct {
	// Cached removes entries from the index while keeping working tree files.
	Cached bool
	// DryRun reports the paths that would be removed without mutating state.
	DryRun bool
	// Recursive permits removing tracked descendants of directory pathspecs.
	Recursive bool
	// Force skips safety checks unless DryRun is also enabled.
	Force bool
}

// RemoveResult returns the tracked paths removed by an rm invocation.
type RemoveResult struct {
	Removed []domain.NormalizedPath
}

// RemoveService removes tracked files from the index and optionally the working tree.
type RemoveService struct {
	indexService   *core.IndexService
	treeResolver   *core.TreeResolver
	changeDetector *core.ChangeDetector
	workspace      *domain.Workspace
}

type removePlan struct {
	paths      []domain.NormalizedPath
	targets    map[domain.NormalizedPath]bool
	pruneRoots []domain.NormalizedPath
}

type fileBackup struct {
	mode os.FileMode
	body []byte
}

// NewRemoveService creates an rm service with the required dependencies.
func NewRemoveService(
	indexService *core.IndexService,
	treeResolver *core.TreeResolver,
	changeDetector *core.ChangeDetector,
	workspace *domain.Workspace,
) *RemoveService {
	return &RemoveService{
		indexService:   indexService,
		treeResolver:   treeResolver,
		changeDetector: changeDetector,
		workspace:      workspace,
	}
}

// Remove removes tracked paths from the index and optionally from the working tree.
func (r *RemoveService) Remove(pathspecs []string, options RemoveOptions) (*RemoveResult, error) {
	index, err := r.indexService.Read()
	if err != nil {
		return nil, fmt.Errorf("rm: %w", err)
	}

	plan, err := r.collectPlan(index, pathspecs, options.Recursive)
	if err != nil {
		return nil, wrapRemoveError(err)
	}

	if !shouldBypassSafetyChecks(options) {
		headPathHashes, err := r.loadHeadPathHashes()
		if err != nil {
			return nil, wrapRemoveError(err)
		}
		if err := r.validateTargets(index, plan.paths, headPathHashes, options.Cached); err != nil {
			return nil, wrapRemoveError(err)
		}
	}

	result := &RemoveResult{
		Removed: append([]domain.NormalizedPath(nil), plan.paths...),
	}
	if options.DryRun {
		return result, nil
	}

	updatedIndex := cloneIndexWithoutTargets(index, plan.targets)
	if options.Cached {
		if err := r.indexService.Write(updatedIndex); err != nil {
			return nil, fmt.Errorf("rm: %w", err)
		}
		return result, nil
	}

	if err := r.applyRemoval(updatedIndex, plan); err != nil {
		return nil, wrapRemoveError(err)
	}
	return result, nil
}

func (r *RemoveService) collectPlan(
	index *domain.Index,
	pathspecs []string,
	recursive bool,
) (removePlan, error) {
	targets := make(map[domain.NormalizedPath]bool)
	pruneRoots := make(map[domain.NormalizedPath]bool)

	for _, pathspec := range pathspecs {
		normPath, err := r.normalizePathspec(pathspec)
		if err != nil {
			return removePlan{}, err
		}

		if entry, _ := index.FindEntry(normPath); entry != nil {
			targets[entry.Path] = true
			continue
		}

		displayPath := removeDisplayPath(pathspec, normPath)
		descendants := findDescendantEntries(index, normPath)
		if len(descendants) == 0 {
			return removePlan{}, newRemovePathDidNotMatchError(displayPath)
		}
		if !recursive {
			return removePlan{}, newRemoveRecursiveRequiredError(displayPath)
		}

		pruneRoots[normPath] = true
		for _, entry := range descendants {
			targets[entry.Path] = true
		}
	}

	return removePlan{
		paths:      domain.SortedPathSet(targets),
		targets:    targets,
		pruneRoots: domain.SortedPathSet(pruneRoots),
	}, nil
}

func (r *RemoveService) normalizePathspec(pathspec string) (domain.NormalizedPath, error) {
	absPath, err := filepath.Abs(filepath.FromSlash(pathspec))
	if err != nil {
		return "", err
	}

	relPath, err := filepath.Rel(r.workspace.RepoDir, absPath)
	if err != nil {
		return "", err
	}

	normPath := filepath.ToSlash(relPath)
	if normPath == ".." || strings.HasPrefix(normPath, "../") {
		return "", newRemoveOutsideRepositoryError(pathspec)
	}

	path, err := domain.NewNormalizedPathUnchecked(normPath)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (r *RemoveService) loadHeadPathHashes() (core.PathHashes, error) {
	headPathHashes, err := r.treeResolver.ResolveHEAD()
	if errors.Is(err, core.ErrRefNotFound) {
		return make(core.PathHashes), nil
	}
	if err != nil {
		return nil, err
	}
	return headPathHashes, nil
}

func (r *RemoveService) validateTargets(
	index *domain.Index,
	paths []domain.NormalizedPath,
	headPathHashes core.PathHashes,
	cached bool,
) error {
	for _, path := range paths {
		entry, _ := index.FindEntry(path)
		if entry == nil {
			continue
		}

		staged := hasStagedChanges(path, entry.Hash, headPathHashes)
		local, err := r.hasLocalModifications(entry, !cached || staged)
		if err != nil {
			return err
		}

		switch {
		case staged && local:
			return newRemoveHasStagedAndLocalStateError(path)
		case staged:
			return newRemoveHasStagedChangesError(path)
		case local:
			return newRemoveHasLocalModificationsError(path)
		}
	}
	return nil
}

func (r *RemoveService) hasLocalModifications(entry *domain.IndexEntry, check bool) (bool, error) {
	if !check {
		return false, nil
	}

	changeResult, err := r.changeDetector.DetectFileChange(entry)
	if err != nil {
		return false, err
	}
	return changeResult.FileState == core.FileStateModified, nil
}

func (r *RemoveService) applyRemoval(updatedIndex *domain.Index, plan removePlan) error {
	backups, err := r.captureFileBackups(plan.paths)
	if err != nil {
		return err
	}

	if err := r.deleteWorkingTreeFiles(plan.paths); err != nil {
		return r.restoreBackups(err, backups)
	}
	if err := r.pruneEmptyDirectories(plan.paths, plan.pruneRoots); err != nil {
		return r.restoreBackups(err, backups)
	}
	if err := r.indexService.Write(updatedIndex); err != nil {
		return r.restoreBackups(err, backups)
	}
	return nil
}

func (r *RemoveService) captureFileBackups(paths []domain.NormalizedPath) (
	map[domain.NormalizedPath]fileBackup, error,
) {
	backups := make(map[domain.NormalizedPath]fileBackup)
	for _, path := range paths {
		absPath, err := path.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return nil, err
		}

		info, err := os.Stat(absPath.String())
		switch {
		case errors.Is(err, os.ErrNotExist):
			continue
		case err != nil:
			return nil, fmt.Errorf("failed to stat '%s': %w", absPath, err)
		}

		body, err := os.ReadFile(absPath.String())
		if err != nil {
			return nil, fmt.Errorf("failed to read '%s': %w", absPath, err)
		}
		backups[path] = fileBackup{
			mode: info.Mode().Perm(),
			body: body,
		}
	}
	return backups, nil
}

func (r *RemoveService) deleteWorkingTreeFiles(paths []domain.NormalizedPath) error {
	for _, path := range paths {
		absolutePath, err := path.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return err
		}

		err = os.Remove(absolutePath.String())
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to delete '%s': %w", absolutePath, err)
		}
	}
	return nil
}

func (r *RemoveService) pruneEmptyDirectories(paths, pruneRoots []domain.NormalizedPath) error {
	for _, pruneRoot := range pruneRoots {
		rootPath, err := pruneRoot.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return err
		}

		for _, path := range paths {
			if !pathWithinRoot(path, pruneRoot) {
				continue
			}

			absolutePath, err := path.ToAbsolutePath(r.workspace.RepoDir)
			if err != nil {
				return err
			}
			if err := pruneEmptyParentDirsWithin(
				absolutePath.String(),
				r.workspace.RepoDir,
				rootPath.String(),
			); err != nil {
				return fmt.Errorf("failed to prune empty directories for '%s': %w", absolutePath, err)
			}
		}
	}
	return nil
}

func (r *RemoveService) restoreBackups(
	removeErr error,
	backups map[domain.NormalizedPath]fileBackup,
) error {
	if err := r.restoreFileBackups(backups); err != nil {
		return fmt.Errorf("%w (rollback failed: %v)", removeErr, err)
	}
	return removeErr
}

func (r *RemoveService) restoreFileBackups(backups map[domain.NormalizedPath]fileBackup) error {
	paths := make(map[domain.NormalizedPath]bool, len(backups))
	for path := range backups {
		paths[path] = true
	}

	for _, path := range domain.SortedPathSet(paths) {
		backup := backups[path]
		absPath, err := path.ToAbsolutePath(r.workspace.RepoDir)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(absPath.String()), domain.DirPermission); err != nil {
			return fmt.Errorf("failed to create directory for '%s': %w", absPath, err)
		}
		if err := os.WriteFile(absPath.String(), backup.body, backup.mode); err != nil {
			return fmt.Errorf("failed to restore '%s': %w", absPath, err)
		}
	}
	return nil
}

func findDescendantEntries(index *domain.Index, path domain.NormalizedPath) []*domain.IndexEntry {
	prefix := path.String()
	if prefix != "" {
		prefix += "/"
	}
	return index.FindEntriesByPathPrefix(prefix)
}

func hasStagedChanges(
	path domain.NormalizedPath,
	indexHash domain.Hash,
	headPathHashes core.PathHashes,
) bool {
	headHash, ok := headPathHashes[path]
	return !ok || headHash != indexHash
}

func cloneIndexWithoutTargets(
	index *domain.Index,
	targets map[domain.NormalizedPath]bool,
) *domain.Index {
	clonedIndex := index.Clone()

	entries := make([]*domain.IndexEntry, 0, len(clonedIndex.Entries))
	for _, entry := range clonedIndex.Entries {
		if !targets[entry.Path] {
			entries = append(entries, entry)
		}
	}

	clonedIndex.ReplaceEntries(entries)
	return clonedIndex
}

func pathWithinRoot(path, root domain.NormalizedPath) bool {
	if root == domain.RootPath {
		return true
	}
	rootPrefix := root.String() + "/"
	return path.String() == root.String() || strings.HasPrefix(path.String(), rootPrefix)
}

func pruneEmptyParentDirsWithin(filePath, repoRoot, pruneRoot string) error {
	dir := filepath.Dir(filePath)
	repoRoot = filepath.Clean(repoRoot)
	pruneRoot = filepath.Clean(pruneRoot)

	stopDir := repoRoot
	if pruneRoot != repoRoot {
		stopDir = filepath.Dir(pruneRoot)
	}

	for {
		dir = filepath.Clean(dir)
		if dir == stopDir || !dirWithinRoot(dir, pruneRoot) {
			return nil
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				dir = filepath.Dir(dir)
				continue
			}
			return err
		}
		if len(entries) != 0 {
			return nil
		}

		if err := os.Remove(dir); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				dir = filepath.Dir(dir)
				continue
			}
			return err
		}
		dir = filepath.Dir(dir)
	}
}

func dirWithinRoot(dir, root string) bool {
	if dir == root {
		return true
	}
	return strings.HasPrefix(dir, root+string(filepath.Separator))
}

func shouldBypassSafetyChecks(options RemoveOptions) bool {
	return options.Force && !options.DryRun
}

func removeDisplayPath(pathspec string, normPath domain.NormalizedPath) string {
	if normPath == domain.RootPath {
		return pathspec
	}
	return normPath.String()
}

func wrapRemoveError(err error) error {
	if errors.Is(err, errRemovePathDidNotMatch) ||
		errors.Is(err, errRemoveRecursiveRequired) ||
		errors.Is(err, errRemoveOutsideRepository) ||
		errors.Is(err, errRemoveHasStagedChanges) ||
		errors.Is(err, errRemoveHasLocalModifications) ||
		errors.Is(err, errRemoveHasStagedAndLocalState) {
		return err
	}
	return fmt.Errorf("rm: %w", err)
}
