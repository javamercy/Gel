package diff

// Region represents a contiguous slice of line diffs included in a hunk.
type Region struct {
	Start int
	End   int
}

// Hunk describes one unified-diff section with old/new ranges and lines.
type Hunk struct {
	OldStart  int
	OldLength int
	NewStart  int
	NewLength int
	Lines     []LineDiff
}

// FindRegions finds changed regions and extends them with fixed context lines.
func FindRegions(diffs []LineDiff) []Region {
	regions := make([]Region, 0)
	size := len(diffs)
	for i := 0; i < size; i++ {
		if diffs[i].OperationType == OpTypeMatch {
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

// BuildHunks converts regions into unified-style hunks with computed ranges.
func BuildHunks(lines []LineDiff, regions []Region) []*Hunk {
	hunks := make([]*Hunk, 0)
	for _, region := range regions {
		oldLength := 0
		newLength := 0

		for i := region.Start; i < region.End; i++ {
			switch lines[i].OperationType {
			case OpTypeMatch:
				oldLength++
				newLength++
			case OpTypeInsertion:
				newLength++
			case OpTypeDeletion:
				oldLength++
			}
		}
		oldStart := lines[region.Start].OldPos
		newStart := lines[region.Start].NewPos
		hunks = append(
			hunks, &Hunk{
				OldStart:  oldStart,
				OldLength: oldLength,
				NewStart:  newStart,
				NewLength: newLength,
				Lines:     lines[region.Start:region.End],
			},
		)
	}
	return hunks
}
