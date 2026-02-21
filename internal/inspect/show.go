package inspect

import (
	"Gel/domain"
	"Gel/internal/core"
	"Gel/internal/diff"
	"Gel/internal/workspace"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type ShowService struct {
	objectService *core.ObjectService
	refService    *core.RefService
	diffService   *diff.DiffService
}

func NewShowService(
	objectService *core.ObjectService,
	refService *core.RefService,
	diffService *diff.DiffService,
) *ShowService {
	return &ShowService{
		objectService: objectService,
		refService:    refService,
		diffService:   diffService,
	}
}

func (s *ShowService) Show(writer io.Writer, objectRef string) error {
	if objectRef == "" {
		return s.showHEAD(writer)
	}

	ref := filepath.Join(workspace.RefsDirName, workspace.HeadsDirName, objectRef)
	if branchExists := s.refService.Exists(ref); branchExists {
		return s.showBranch(writer, ref)
	}

	object, err := s.objectService.Read(objectRef)
	if err != nil {
		return err
	}
	switch object.Type() {
	case domain.ObjectTypeBlob:
		return s.showBlob(writer, object.(*domain.Blob))
	case domain.ObjectTypeTree:

		return s.showTree(writer, object.(*domain.Tree), objectRef)
	case domain.ObjectTypeCommit:
		headRef, err := s.refService.ReadSymbolic(workspace.HeadFileName)
		if err != nil {
			return err
		}
		return s.showCommit(writer, objectRef, s.trimBranchName(headRef))
	default:
		return fmt.Errorf("unsupported object type: %s", object.Type())
	}
}

func (s *ShowService) showHEAD(writer io.Writer) error {
	commitHash, err := s.refService.Resolve(workspace.HeadFileName)
	if err != nil {
		return err
	}
	ref, err := s.refService.ReadSymbolic(workspace.HeadFileName)
	if err != nil {
		return err
	}

	return s.showCommit(writer, commitHash, s.trimBranchName(ref))
}

func (s *ShowService) showBranch(writer io.Writer, ref string) error {
	commitHash, err := s.refService.Read(ref)
	if err != nil {
		return err
	}
	return s.showCommit(writer, commitHash, s.trimBranchName(ref))
}

func (s *ShowService) showCommit(writer io.Writer, commitHash string, branchName string) error {
	commit, err := s.objectService.ReadCommit(commitHash)
	if err != nil {
		return err
	}

	parentCommitHash := ""
	if len(commit.ParentHashes) > 0 {
		parentCommitHash = commit.ParentHashes[0]
	}

	options := diff.DiffOptions{
		Mode:             diff.ModeCommitVsCommit,
		BaseCommitHash:   parentCommitHash,
		TargetCommitHash: commitHash,
	}
	commitDate, err := domain.FormatCommitDate(
		commit.Author.Timestamp,
		commit.Author.Timezone,
	)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		writer,
		"commit %s (HEAD -> %s)\n"+
			"Author: %s <%s>\n"+
			"Date: %s\n\n\t%s\n\n",
		commitHash, branchName,
		commit.Author.Name, commit.Author.Email,
		commitDate, commit.Message,
	); err != nil {
		return err
	}
	return s.diffService.Diff(writer, options)
}

func (s *ShowService) showBlob(writer io.Writer, blob *domain.Blob) error {
	_, err := fmt.Fprintf(writer, "%s", blob.Body())
	return err
}

func (s *ShowService) showTree(writer io.Writer, tree *domain.Tree, treeHash string) error {
	entries, err := tree.Deserialize()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "tree %s\n\n", treeHash); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%s\n", entry.Name); err != nil {
			return err
		}
	}
	return nil
}

func (s *ShowService) trimBranchName(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}
