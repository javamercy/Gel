package diff

import (
	"Gel/domain"
	"Gel/internal/core"
	"fmt"
	"io"
	"os"
	"strings"
)

type DiffMode int

const (
	ModeWorkingTreeVsIndex  DiffMode = iota // default: gel diff
	ModeWorkingTreeVsHEAD                   // gel diff HEAD
	ModeIndexVsHEAD                         // gel diff --staged
	ModeCommitVsWorkingTree                 // gel diff <commit>
	ModeCommitVsCommit                      // gel diff <commit1> <commit2>
)

type DiffOptions struct {
	Mode             DiffMode
	BaseCommitHash   string
	TargetCommitHash string
}

type ContentLoaderFunc func(path, hash string) (string, error)

type DiffStatus int

const (
	DiffStatusModified DiffStatus = iota
	DiffStatusAdded
	DiffStatusDeleted
)

type DiffResult struct {
	Hunks   []Hunk
	Status  DiffStatus
	OldPath string
	NewPath string
	OldHash string
	NewHash string
}
type DiffService struct {
	objectService *core.ObjectService
	refService    *core.RefService
	treeResolver  *core.TreeResolver
	diffAlgorithm *MyersDiffAlgorithm
}

func NewDiffService(
	objectService *core.ObjectService,
	refService *core.RefService,
	treeResolver *core.TreeResolver,
	diffAlgorithm *MyersDiffAlgorithm,
) *DiffService {
	return &DiffService{
		objectService: objectService,
		refService:    refService,
		treeResolver:  treeResolver,
		diffAlgorithm: diffAlgorithm,
	}
}

func (d *DiffService) Diff(writer io.Writer, options DiffOptions) error {
	var results []*DiffResult
	var resultsErr error
	switch options.Mode {
	case ModeWorkingTreeVsHEAD:
		headEntries, err := d.treeResolver.ResolveHEAD()
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		results, resultsErr = d.ComputeDiffResults(
			headEntries, workingTreeEntries, d.LoadBlobContent, d.LoadFileContent,
		)
	case ModeIndexVsHEAD:
		headEntries, err := d.treeResolver.ResolveHEAD()
		if err != nil {
			return err
		}
		indexEntries, err := d.treeResolver.ResolveIndex()
		if err != nil {
			return err
		}
		results, resultsErr = d.ComputeDiffResults(
			indexEntries, headEntries, d.LoadBlobContent, d.LoadBlobContent,
		)
	case ModeCommitVsCommit:
		var baseCommitEntries map[string]string
		if options.BaseCommitHash == "" {
			baseCommitEntries = make(map[string]string)
		} else {
			var err error
			baseCommitEntries, err = d.treeResolver.ResolveCommit(options.BaseCommitHash)
			if err != nil {
				return err
			}
		}

		targetCommitEntries, err := d.treeResolver.ResolveCommit(options.TargetCommitHash)
		if err != nil {
			return err
		}
		results, resultsErr = d.ComputeDiffResults(
			targetCommitEntries, baseCommitEntries, d.LoadBlobContent, d.LoadBlobContent,
		)
	case ModeCommitVsWorkingTree:
		commitEntries, err := d.treeResolver.ResolveCommit(options.BaseCommitHash)
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		results, resultsErr = d.ComputeDiffResults(
			workingTreeEntries, commitEntries, d.LoadFileContent, d.LoadBlobContent,
		)
	case ModeWorkingTreeVsIndex:
		indexEntries, err := d.treeResolver.ResolveIndex()
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		results, resultsErr = d.ComputeDiffResults(
			workingTreeEntries, indexEntries, d.LoadFileContent, d.LoadBlobContent,
		)
	default:
		return fmt.Errorf("unsupported diff mode: %d", options.Mode)
	}
	if resultsErr != nil {
		return resultsErr
	}
	return d.PrintResults(writer, results)
}

func (d *DiffService) ComputeDiffResults(
	newEntries, oldEntries map[string]string,
	newContentLoader, oldContentLoader ContentLoaderFunc,
) ([]*DiffResult, error) {
	var results []*DiffResult
	for newPath, newHash := range newEntries {
		newData, err := newContentLoader(newPath, newHash)
		if err != nil {
			return nil, err
		}

		newLines := strings.Split(strings.TrimSuffix(newData, "\n"), "\n")
		oldHash, ok := oldEntries[newPath]
		if !ok {
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs([]string{}, newLines)
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			results = append(
				results, &DiffResult{
					hunks, DiffStatusAdded, "",
					newPath, "", newHash,
				},
			)
		} else if newHash != oldHash {
			oldData, err := oldContentLoader("", oldHash)
			if err != nil {
				return nil, err
			}

			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, newLines)
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			results = append(
				results, &DiffResult{
					hunks, DiffStatusModified, "",
					newPath, oldHash, newHash,
				},
			)
		}
	}

	for oldPath, oldHash := range oldEntries {
		if _, ok := newEntries[oldPath]; !ok {
			oldData, err := oldContentLoader(oldPath, oldHash)
			if err != nil {
				return nil, err
			}

			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, []string{})
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			results = append(
				results, &DiffResult{
					hunks, DiffStatusDeleted, oldPath,
					"", oldHash, "",
				},
			)
		}
	}
	return results, nil
}

func (d *DiffService) LoadBlobContent(_, hash string) (string, error) {
	blob, err := d.objectService.ReadBlob(hash)
	if err != nil {
		return "", err
	}
	return string(blob.Body()), nil
}

func (d *DiffService) LoadFileContent(path, _ string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *DiffService) PrintResults(writer io.Writer, results []*DiffResult) error {
	for _, result := range results {
		var err error
		switch result.Status {
		case DiffStatusAdded:
			err = d.PrintAddedFileHeader(writer, result.OldPath, result.NewPath, result.NewHash[:7])

		case DiffStatusModified:
			err = d.PrintModifiedFileHeader(
				writer, result.OldPath, result.NewPath, result.OldHash[:7], result.NewHash[:7],
			)
		case DiffStatusDeleted:
			err = d.PrintDeletedFileHeader(writer, result.OldPath, result.OldHash[:7])
		}
		if err != nil {
			return err
		}
		if err := d.PrintHunks(writer, result.Hunks); err != nil {
			return err
		}
	}
	return nil
}

func (d *DiffService) PrintAddedFileHeader(writer io.Writer, oldPath, newPath, hash string) error {
	if _, err := fmt.Fprintf(
		writer, "%sdiff --gel a/%s b/%s%s\n", core.ColorBold, oldPath, newPath, core.ColorReset,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		writer, "%snew file mode %s%s\n", core.ColorBold, domain.RegularFileMode, core.ColorReset,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%sindex 00000000..%s%s\n", core.ColorBold, hash, core.ColorReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s--- /dev/null%s\n", core.ColorBold, core.ColorReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s+++ b/%s%s\n", core.ColorBold, newPath, core.ColorReset); err != nil {
		return err
	}
	return nil
}

func (d *DiffService) PrintDeletedFileHeader(writer io.Writer, oldPath, oldHash string) error {
	if _, err := fmt.Fprintf(
		writer, "%sdeleted file mode %s%s\n", core.ColorBold, domain.RegularFileMode, core.ColorReset,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(
		writer, "%sindex %s..00000000%s\n", core.ColorBold, oldHash, core.ColorReset,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s--- a/%s%s\n", core.ColorBold, oldPath, core.ColorReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s+++ /dev/null%s\n", core.ColorBold, core.ColorReset); err != nil {
		return err
	}
	return nil
}

func (d *DiffService) PrintModifiedFileHeader(writer io.Writer, oldPath, newPath, oldHash, newHash string) error {
	if _, err := fmt.Fprintf(
		writer, "%sindex %s..%s %s%s\n", core.ColorBold, oldHash, newHash, domain.RegularFileMode, core.ColorReset,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s--- a/%s%s\n", core.ColorBold, oldPath, core.ColorReset); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "%s+++ b/%s%s\n", core.ColorBold, newPath, core.ColorReset); err != nil {
		return err
	}
	return nil
}

func (d *DiffService) PrintHunks(writer io.Writer, hunks []Hunk) error {
	for _, hunk := range hunks {
		header := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunk.OldStart, hunk.OldLength, hunk.NewStart, hunk.NewLength)
		if _, err := fmt.Fprintln(writer, header); err != nil {
			return err
		}
		for _, line := range hunk.Lines {
			var prefix string
			var color string
			switch line.OperationType {
			case OpTypeMatch:
				prefix = " "
				color = ""
			case OpTypeInsertion:
				prefix = "+ "
				color = core.ColorGreen
			case OpTypeDeletion:
				prefix = "- "
				color = core.ColorRed
			}
			if _, err := fmt.Fprintf(writer, "%s%s%s%s\n", color, prefix, line.Content, core.ColorReset); err != nil {
				return err
			}
		}
	}
	return nil
}
