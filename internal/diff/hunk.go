package diff

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

func BuildHunks(diffs []LineDiff, regions []Region) []Hunk {
	hunks := make([]Hunk, 0)
	for _, region := range regions {
		oldLength := 0
		newLength := 0

		for i := region.Start; i < region.End; i++ {
			switch diffs[i].OperationType {
			case OpTypeMatch:
				oldLength++
				newLength++
			case OpTypeInsertion:
				newLength++
			case OpTypeDeletion:
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
