package diff

type OperationType int

const (
	OpTypeInsertion OperationType = iota
	OpTypeDeletion
	OpTypeMatch
)

func (o OperationType) String() string {
	switch o {
	case OpTypeInsertion:
		return "insertion"
	case OpTypeDeletion:
		return "deletion"
	case OpTypeMatch:
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
