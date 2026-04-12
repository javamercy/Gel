package domain

type ObjectType string

const (
	ObjectTypeBlob   ObjectType = "blob"
	ObjectTypeTree   ObjectType = "tree"
	ObjectTypeCommit ObjectType = "commit"
)

func (objectType ObjectType) IsValid() bool {
	switch objectType {
	case ObjectTypeBlob, ObjectTypeTree, ObjectTypeCommit:
		return true
	default:
		return false
	}
}

func (objectType ObjectType) String() string {
	return string(objectType)
}

func ParseObjectType(typeStr string) (ObjectType, bool) {
	objectType := ObjectType(typeStr)
	if objectType.IsValid() {
		return objectType, true
	}
	return "", false
}
