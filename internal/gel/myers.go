package gel

import (
	"slices"
)

type OperationType int

const (
	Insertion OperationType = iota
	Deletion
	Match
)

func (o OperationType) String() string {
	switch o {
	case Insertion:
		return "insertion"
	case Deletion:
		return "deletion"
	case Match:
		return "match"
	}
	return ""
}

type LineDiff struct {
	OperationType OperationType
	Content       string
	OldPos        int
	NewPos        int
}

type MyersDiffAlgorithm struct {
}

func NewMyersDiffAlgorithm() *MyersDiffAlgorithm {
	return &MyersDiffAlgorithm{}
}

func (m *MyersDiffAlgorithm) ComputeLineDiffs(A, B []string) []LineDiff {
	trace := m.computeTrace(A, B)
	return m.backtrack(trace, A, B)
}

func (m *MyersDiffAlgorithm) computeTrace(A, B []string) [][]int {
	lenA, lenB := len(A), len(B)
	if lenA == 0 && lenB == 0 {
		return nil
	}
	maxDepth := lenA + lenB
	trace := make([][]int, 0)
	v := make([]int, 2*maxDepth+1)
	for d := 0; d <= maxDepth; d++ {
		vCopy := make([]int, len(v))
		for k := -d; k <= d; k += 2 {
			var x, y int
			if k == -d || (k != d && v[maxDepth+k-1] < v[maxDepth+k+1]) {
				x = v[maxDepth+k+1]
			} else {
				x = v[maxDepth+k-1] + 1
			}

			y = x - k
			for x < lenA && y < lenB && A[x] == B[y] {
				x++
				y++
			}
			v[maxDepth+k] = x
			if x == lenA && y == lenB {
				copy(vCopy, v)
				trace = append(trace, vCopy)
				return trace
			}
		}
		copy(vCopy, v)
		trace = append(trace, vCopy)
	}
	return trace
}

func (m *MyersDiffAlgorithm) backtrack(trace [][]int, A, B []string) []LineDiff {
	diffs := make([]LineDiff, 0)
	lenA, lenB := len(A), len(B)
	x, y := lenA, lenB
	maxDepth := lenA + lenB
	for d := len(trace) - 1; d > 0; d-- {
		k := x - y
		v := trace[d-1]
		var prevK int
		if k == -d || (k != d && v[maxDepth+k-1] < v[maxDepth+k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX := v[maxDepth+prevK]
		prevY := prevX - prevK
		for x > prevX && y > prevY {
			diffs = append(
				diffs,
				LineDiff{Match, A[x-1], x, y},
			)
			x--
			y--
		}
		if x > prevX {
			diffs = append(
				diffs,
				LineDiff{Deletion, A[x-1], x, y},
			)
			x--
		} else {
			diffs = append(
				diffs,
				LineDiff{Insertion, B[y-1], x, y},
			)
			y--
		}
	}

	if x > 0 && y > 0 {
		for x > 0 && y > 0 {
			diffs = append(diffs, LineDiff{Match, A[x-1], x, y})
			x--
			y--
		}
	} else if x > 0 && y == 0 {
		for x > 0 {
			diffs = append(diffs, LineDiff{Deletion, A[x-1], x, y})
			x--
		}
	} else if y > 0 {
		for y > 0 {
			diffs = append(diffs, LineDiff{Match, A[y-1], x, y})
			y--
		}
	}
	slices.Reverse(diffs)
	return diffs
}
