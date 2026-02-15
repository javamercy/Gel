package gel

import (
	"Gel/domain"
	"Gel/internal/gel/diff"
	"fmt"
	"os"
	"strings"
)

type ContentLoaderFunc func(path, hash string) (string, error)
type DiffService struct {
	objectService *ObjectService
	refService    *RefService
	treeResolver  *TreeResolver
	diffAlgorithm *diff.MyersDiffAlgorithm
}

func NewDiffService(
	objectService *ObjectService,
	refService *RefService,
	treeResolver *TreeResolver,
	diffAlgorithm *diff.MyersDiffAlgorithm,
) *DiffService {
	return &DiffService{
		objectService: objectService,
		refService:    refService,
		treeResolver:  treeResolver,
		diffAlgorithm: diffAlgorithm,
	}
}

func (d *DiffService) Diff(
	head, staged bool,
	baseCommitHash, targetCommitHash string,
) error {
	if head {
		headEntries, err := d.treeResolver.ResolveHEAD()
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(headEntries, workingTreeEntries, d.loadBlobContent, d.loadFileContent)
	} else if staged {
		headEntries, err := d.treeResolver.ResolveHEAD()
		if err != nil {
			return err
		}
		indexEntries, err := d.treeResolver.ResolveIndex()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(indexEntries, headEntries, d.loadBlobContent, d.loadBlobContent)
	} else if baseCommitHash != "" && targetCommitHash != "" {
		baseCommitEntries, err := d.treeResolver.ResolveCommit(baseCommitHash)
		if err != nil {
			return err
		}
		targetCommitEntries, err := d.treeResolver.ResolveCommit(targetCommitHash)
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(baseCommitEntries, targetCommitEntries, d.loadBlobContent, d.loadBlobContent)
	} else if baseCommitHash != "" {
		commitEntries, err := d.treeResolver.ResolveCommit(baseCommitHash)
		if err != nil {
			return err
		}
		workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
		if err != nil {
			return err
		}
		return d.computeEntryDiffs(workingTreeEntries, commitEntries, d.loadFileContent, d.loadBlobContent)
	}

	indexEntries, err := d.treeResolver.ResolveIndex()
	if err != nil {
		return err
	}
	workingTreeEntries, err := d.treeResolver.ResolveWorkingTree()
	if err != nil {
		return err
	}
	return d.computeEntryDiffs(workingTreeEntries, indexEntries, d.loadFileContent, d.loadBlobContent)
}

func (d *DiffService) computeEntryDiffs(
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
			d.printNewFileHeader(newPath, newPath, newHash[:7])
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs([]string{}, newLines)
			regions := diff.FindRegions(lineDiffs)
			hunks := diff.BuildHunks(lineDiffs, regions)
			d.printHunks(hunks)
		} else if newHash != oldHash {
			oldData, err := oldContentLoader("", oldHash)
			if err != nil {
				return err
			}
			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, newLines)
			regions := diff.FindRegions(lineDiffs)
			hunks := diff.BuildHunks(lineDiffs, regions)
			d.printModifiedFileHeader(newPath, newPath, oldHash[:7], newHash[:7])
			d.printHunks(hunks)
		}
	}

	for oldPath, oldHash := range oldEntries {
		if _, ok := newEntries[oldPath]; !ok {
			d.printDeletedFileHeader(oldPath, oldHash[:7])
			oldData, err := oldContentLoader(oldPath, oldHash)
			if err != nil {
				return err
			}

			oldLines := strings.Split(strings.TrimSuffix(oldData, "\n"), "\n")
			lineDiffs := d.diffAlgorithm.ComputeLineDiffs(oldLines, []string{})
			regions := diff.FindRegions(lineDiffs)
			hunks := diff.BuildHunks(lineDiffs, regions)
			d.printHunks(hunks)
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

func (d *DiffService) printNewFileHeader(oldPath, newPath, hash string) {
	fmt.Printf("%sdiff --gel a/%s b/%s%s\n", colorBold, oldPath, newPath, colorReset)
	fmt.Printf("%snew file mode %s%s\n", colorBold, domain.RegularFileMode, colorReset)
	fmt.Printf("%sindex 00000000..%s%s\n", colorBold, hash, colorReset)
	fmt.Printf("%s--- /dev/null%s\n", colorBold, colorReset)
	fmt.Printf("%s+++ b/%s%s\n", colorBold, newPath, colorReset)
}

func (d *DiffService) printDeletedFileHeader(oldPath, oldHash string) {
	fmt.Printf("%sdeleted file mode %s%s\n", colorBold, domain.RegularFileMode, colorReset)
	fmt.Printf("%sindex %s..00000000%s\n", colorBold, oldHash, colorReset)
	fmt.Printf("%s--- a/%s%s\n", colorBold, oldPath, colorReset)
	fmt.Printf("%s+++ /dev/null%s\n", colorBold, colorReset)
}

func (d *DiffService) printModifiedFileHeader(oldPath, newPath, oldHash, newHash string) {
	fmt.Printf("%sindex %s..%s %s%s\n", colorBold, oldHash, newHash, domain.RegularFileMode, colorReset)
	fmt.Printf("%s--- a/%s%s\n", colorBold, oldPath, colorReset)
	fmt.Printf("%s+++ b/%s%s\n", colorBold, newPath, colorReset)
}

func (d *DiffService) printHunks(hunks []diff.Hunk) {
	for _, hunk := range hunks {
		header := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunk.OldStart, hunk.OldLength, hunk.NewStart, hunk.NewLength)
		fmt.Println(header)
		for _, line := range hunk.Lines {
			var prefix string
			var color string
			switch line.OperationType {
			case diff.OpTypeMatch:
				prefix = " "
				color = ""
			case diff.OpTypeInsertion:
				prefix = "+ "
				color = colorGreen
			case diff.OpTypeDeletion:
				prefix = "- "
				color = colorRed
			}
			fmt.Printf("%s%s%s%s\n", color, prefix, line.Content, colorReset)
		}
	}
}
