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
		return d.computeEntryDiffs(writer, headEntries, workingTreeEntries, d.loadBlobContent, d.loadFileContent)
	case ModeIndexVsHEAD:
		headEntries, err := d.treeResolver.ResolveHEAD()
		if err != nil {
			return err
		}
		indexEntries, err := d.treeResolver.ResolveIndex()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(writer, indexEntries, headEntries, d.loadBlobContent, d.loadBlobContent)
	case ModeCommitVsCommit:
		baseCommitEntries, err := d.treeResolver.ResolveCommit(options.BaseCommitHash)
		if err != nil {
			return err
		}
		targetCommitEntries, err := d.treeResolver.ResolveCommit(options.TargetCommitHash)
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(writer, targetCommitEntries, baseCommitEntries, d.loadBlobContent, d.loadBlobContent)
	case ModeCommitVsWorkingTree:
		commitEntries, err := d.treeResolver.ResolveCommit(options.BaseCommitHash)
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(writer, workingTreeEntries, commitEntries, d.loadFileContent, d.loadBlobContent)
	case ModeWorkingTreeVsIndex:
		indexEntries, err := d.treeResolver.ResolveIndex()
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(writer, workingTreeEntries, indexEntries, d.loadFileContent, d.loadBlobContent)
	default:
		return fmt.Errorf("unsupported diff mode: %d", options.Mode)
	}
}

// TODO: update this func such that it does not print, returns a diff object.
func (d *DiffService) computeEntryDiffs(
	writer io.Writer,
	newEntries, oldEntries map[string]string,
	newContentLoader, oldContentLoader ContentLoaderFunc,
) error {
	for newPath, newHash := range newEntries {
		newData, err := newContentLoader(newPath, newHash)
		if err != nil {
			return err
		}

		newLines := strings.Split(strings.TrimSuffix(newData, "\n"), "\n")
		oldHash, ok := oldEntries[newPath]
		if !ok {
			if err := d.printNewFileHeader(writer, newPath, newPath, newHash[:7]); err != nil {
				return err
			}
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs([]string{}, newLines)
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			if err := d.printHunks(writer, hunks); err != nil {
				return err
			}
		} else if newHash != oldHash {
			oldData, err := oldContentLoader("", oldHash)
			if err != nil {
				return err
			}
			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, newLines)
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			if err := d.printModifiedFileHeader(writer, newPath, newPath, oldHash[:7], newHash[:7]); err != nil {
				return err
			}
			if err := d.printHunks(writer, hunks); err != nil {
				return err
			}
		}
	}

	for oldPath, oldHash := range oldEntries {
		if _, ok := newEntries[oldPath]; !ok {
			if err := d.printDeletedFileHeader(writer, oldPath, oldHash[:7]); err != nil {
				return err
			}
			oldData, err := oldContentLoader(oldPath, oldHash)
			if err != nil {
				return err
			}

			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, []string{})
			regions := FindRegions(lineDiffs)
			hunks := BuildHunks(lineDiffs, regions)
			if err := d.printHunks(writer, hunks); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DiffService) loadBlobContent(_, hash string) (string, error) {
	blob, err := d.objectService.ReadBlob(hash)
	if err != nil {
		return "", err
	}
	return string(blob.Body()), nil
}

func (d *DiffService) loadFileContent(path, _ string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (d *DiffService) printNewFileHeader(writer io.Writer, oldPath, newPath, hash string) error {
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

func (d *DiffService) printDeletedFileHeader(writer io.Writer, oldPath, oldHash string) error {
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

func (d *DiffService) printModifiedFileHeader(writer io.Writer, oldPath, newPath, oldHash, newHash string) error {
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

func (d *DiffService) printHunks(writer io.Writer, hunks []Hunk) error {
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
