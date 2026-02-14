package gel

import (
	"Gel/domain"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Region struct {
	Start int
	End   int
}

type Hunk struct {
	OldStart  int
	OldLength int
	NewStart  int
	NewLength int
	Lines     []LineDiff
}

type DiffService struct {
	objectService     *ObjectService
	indexService      *IndexService
	refService        *RefService
	workingDirService *WorkingDirService
	diffHelper        *MyersDiffAlgorithm
}

func NewDiffService(
	objectService *ObjectService,
	indexService *IndexService,
	refService *RefService,
	workingDirService *WorkingDirService,
	diffHelper *MyersDiffAlgorithm,
) *DiffService {
	return &DiffService{
		objectService:     objectService,
		indexService:      indexService,
		refService:        refService,
		workingDirService: workingDirService,
		diffHelper:        diffHelper,
	}
}

func (d *DiffService) Diff(head, staged bool, firstCommitHash, secondCommitHash string) error {
	if head {
		return d.diffWithHead()
	}

	if staged {
		return d.diffWithStaged(firstCommitHash)
	}

	if firstCommitHash != "" && secondCommitHash != "" {
		return d.diffWithCommits(firstCommitHash, secondCommitHash)
	}
	if firstCommitHash != "" {
		return d.diffWithCommit(firstCommitHash)
	}

	indexEntries, err := d.indexService.GetEntryMap()
	if err != nil {
		return err
	}
	workingDirFiles, err := d.workingDirService.GetFileMap()
	if err != nil {
		return err
	}

	for path, hash := range workingDirFiles {
		workingDirFileData, err := os.ReadFile(path)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			return err
		}

		workingDirFileLines := strings.Split(strings.TrimSuffix(string(workingDirFileData), "\n"), "\n")
		indexHash, ok := indexEntries[path]
		if !ok {
			fmt.Printf("diff --gel a/%s b/%s\n", path, path)
			fmt.Printf("new file mode %s\n", domain.RegularFileStr)
			fmt.Printf("index 0000000..%s %s\n", hash[:7], domain.RegularFileStr)
			fmt.Printf("--- /dev/null\n+++ b/%s\n", path)

			lineDiffs := d.diffHelper.ComputeLineDiffs([]string{}, workingDirFileLines)
			regions := d.FindRegions(lineDiffs)
			hunks := d.BuildHunks(lineDiffs, regions)
			d.printDiff(hunks)
		} else if hash != indexHash {
			blob, err := d.objectService.ReadBlob(indexHash)
			if err != nil {
				return err
			}

			indexFileLines := strings.Split(strings.TrimSuffix(string(blob.Body()), "\n"), "\n")
			lineDiffs := d.diffHelper.ComputeLineDiffs(indexFileLines, workingDirFileLines)
			regions := d.FindRegions(lineDiffs)
			hunks := d.BuildHunks(lineDiffs, regions)
			fmt.Printf("diff --gel a/%s b/%s\n", path, path)
			fmt.Printf("index %s...%s %s\n", hash[:7], indexHash[:7], domain.RegularFileStr)
			fmt.Printf("--- a/%s\n+++ b/%s\n", path, path)
			d.printDiff(hunks)
		} else {
			//fmt.Println("Unchanged file: ", path)
		}
	}

	for path, hash := range indexEntries {
		if _, ok := workingDirFiles[path]; !ok {
			fmt.Printf("diff --gel a/%s b/%s\n", path, path)
			fmt.Printf("deleted file mode %s\n", domain.RegularFileStr)
			fmt.Printf("index %s..0000000 %s\n", hash[:7], domain.RegularFileStr)
			fmt.Printf("--a/%s\n", path)

			indexFileData, err := d.objectService.ReadBlob(hash)
			if err != nil {
				return err
			}

			indexFileLines := strings.Split(strings.TrimSuffix(string(indexFileData.Body()), "\n"), "\n")
			lineDiffs := d.diffHelper.ComputeLineDiffs(indexFileLines, []string{})
			regions := d.FindRegions(lineDiffs)
			hunks := d.BuildHunks(lineDiffs, regions)
			d.printDiff(hunks)
		}
	}
	return nil
}

func (d *DiffService) diffWithStaged(commitHash string) error {
	return nil
}

func (d *DiffService) diffWithHead() error {
	return nil
}

func (d *DiffService) diffWithCommit(commitHash string) error {
	return nil
}

func (d *DiffService) diffWithCommits(firstCommitHash, secondCommitHash string) error {
	return nil
}

func (d *DiffService) FindRegions(diffs []LineDiff) []Region {
	regions := make([]Region, 0)
	size := len(diffs)
	for i := 0; i < size; i++ {
		if diffs[i].OperationType == Match {
			continue
		}
		start, end := i-3, i+4
		if len(regions) == 0 || regions[len(regions)-1].End < start {
			regions = append(regions, Region{max(0, start), min(size, end)})
		} else {
			regions[len(regions)-1].End = min(size, max(regions[len(regions)-1].End, end))
		}
	}
	return regions
}

func (d *DiffService) BuildHunks(diffs []LineDiff, regions []Region) []Hunk {
	hunks := make([]Hunk, 0)
	for _, region := range regions {
		oldLength := 0
		newLength := 0

		for i := region.Start; i < region.End; i++ {
			switch diffs[i].OperationType {
			case Match:
				oldLength++
				newLength++
			case Insertion:
				newLength++
			case Deletion:
				oldLength++
			}
		}
		oldStart := diffs[region.Start].OldPos
		newStart := diffs[region.Start].NewPos
		hunks = append(
			hunks, Hunk{
				OldStart:  oldStart,
				OldLength: oldLength,
				NewStart:  newStart,
				NewLength: newLength,
				Lines:     diffs[region.Start:region.End],
			},
		)
	}
	return hunks
}

func (d *DiffService) printDiff(hunks []Hunk) {
	for _, hunk := range hunks {
		header := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunk.OldStart, hunk.OldLength, hunk.NewStart, hunk.NewLength)
		fmt.Println(header)
		for _, line := range hunk.Lines {
			var prefix string
			var color string
			switch line.OperationType {
			case Match:
				prefix = " "
				color = ""
			case Insertion:
				prefix = "+ "
				color = colorGreen
			case Deletion:
				prefix = "- "
				color = colorRed
			}
			fmt.Printf("%s%s%s%s\n", color, prefix, line.Content, colorReset)
		}
	}
}
