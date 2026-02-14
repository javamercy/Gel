package gel

import (
	"Gel/domain"
	"Gel/internal/workspace"
	"errors"
	"fmt"
	"io"
)

type FileStatus struct {
	Path   string
	Status string
}

type StatusResult struct {
	Staged    []FileStatus
	Unstaged  []FileStatus
	Untracked []string
}
type StatusService struct {
	indexService       *IndexService
	objectService      *ObjectService
	workingDirService  *WorkingDirService
	refService         *RefService
	symbolicRefService *SymbolicRefService
}

func NewStatusService(
	indexService *IndexService, objectService *ObjectService, workingDirService *WorkingDirService,
	refService *RefService, symbolicRefService *SymbolicRefService,
) *StatusService {
	return &StatusService{
		indexService:       indexService,
		objectService:      objectService,
		workingDirService:  workingDirService,
		refService:         refService,
		symbolicRefService: symbolicRefService,
	}
}

func (s *StatusService) Status(writer io.Writer) error {
	result := &StatusResult{}

	indexEntries := make(map[string]string)
	idxEntries, err := s.indexService.GetEntries()
	if err != nil {
		return err
	}
	if len(idxEntries) > 0 {
		for _, entry := range idxEntries {
			indexEntries[entry.Path] = entry.Hash
		}
	}
	headTreeEntries := make(map[string]string)
	commitHash, err := s.refService.Resolve(workspace.HeadFileName)
	if err != nil && !errors.Is(err, ErrRefNotFound) {
		return err
	}
	if commitHash != "" {
		commit, err := s.objectService.ReadCommit(commitHash)
		if err != nil {
			return err
		}
		walkOptions := WalkOptions{
			Recursive:    true,
			IncludeTrees: false,
			OnlyTrees:    false,
		}
		processor := func(entry domain.TreeEntry, relPath string) error {
			headTreeEntries[relPath] = entry.Hash
			return nil
		}
		treeWalker := NewTreeWalker(s.objectService, walkOptions)
		if err := treeWalker.Walk(commit.TreeHash, "", processor); err != nil {
			return err
		}
	}

	workingDirFiles, err := s.workingDirService.GetFileMap()
	if err != nil {
		return err
	}

	// Compare HEAD vs Index → Staged changes
	for indexEntryPath, indexEntryHash := range indexEntries {
		headHash, inHead := headTreeEntries[indexEntryPath]
		if !inHead {
			// in Index but not in HEAD
			result.Staged = append(result.Staged, FileStatus{indexEntryPath, "New File"})
		} else if headHash != indexEntryHash {
			// in Index and in HEAD but different
			result.Staged = append(result.Staged, FileStatus{indexEntryPath, "Modified"})
		}
	}

	for path := range headTreeEntries {
		if _, inIndex := indexEntries[path]; !inIndex {
			// in HEAD but not in Index
			result.Unstaged = append(result.Staged, FileStatus{path, "Deleted"})
		}
	}

	// Compare Index vs Working Dir → Unstaged changes
	for indexEntryPath, indexEntryHash := range indexEntries {
		workingDirHash, inWorkingDir := workingDirFiles[indexEntryPath]
		if !inWorkingDir {
			// in Index but not in Working Dir
			result.Unstaged = append(result.Unstaged, FileStatus{indexEntryPath, "Deleted"})
		} else if workingDirHash != indexEntryHash {
			// in Index and Working Dir but different
			result.Unstaged = append(result.Unstaged, FileStatus{indexEntryPath, "Modified"})
		}
	}

	// Find untracked files
	for path := range workingDirFiles {
		if _, inIndex := indexEntries[path]; !inIndex {
			result.Untracked = append(result.Untracked, path)
		}
	}

	currentBranch, err := s.symbolicRefService.Read(workspace.HeadFileName, true)
	if err != nil {
		return err
	}
	return s.printStatus(writer, currentBranch, len(headTreeEntries), result)
}

func (s *StatusService) printStatus(writer io.Writer, branch string, headTreeSize int, result *StatusResult) error {
	if _, err := fmt.Fprintf(writer, "On branch %s%s%s", colorGreen, branch, colorReset); err != nil {
		return err
	}
	if headTreeSize == 0 {
		if _, err := fmt.Fprintln(writer, " (no commits yet)"); err != nil {
			return err
		}
	}
	if len(result.Staged) > 0 {
		if _, err := fmt.Fprintf(
			writer, "\n%sChanges to be committed:%s\n", colorGreen, colorReset,
		); err != nil {
			return err
		}
		for _, staged := range result.Staged {
			if _, err := fmt.Fprintf(
				writer,
				"\t%s%s:  %s%s\n", colorGreen, staged.Status, staged.Path, colorReset,
			); err != nil {
				return err
			}
		}
	}
	if len(result.Unstaged) > 0 {
		if _, err := fmt.Fprintf(
			writer, "\nChanges not staged for commit:%s\n", colorReset,
		); err != nil {
			return err
		}
		for _, unstaged := range result.Unstaged {
			if _, err := fmt.Fprintf(
				writer,
				"\t%s:  %s%s\n", unstaged.Status, unstaged.Path, colorReset,
			); err != nil {
				return err
			}
		}
	}
	if len(result.Untracked) > 0 {
		if _, err := fmt.Fprintf(writer, "\nUntracked files:%s\n", colorReset); err != nil {
			return err
		}
		for _, untracked := range result.Untracked {
			if _, err := fmt.Fprintf(
				writer,
				"\t%s%s\n", untracked, colorReset,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
