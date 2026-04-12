package domain

// ObjectType represents a type of object in the version control model, such as "blob", "tree", or "commit".
type ObjectType string

const (
	// ObjectTypeBlob represents a file content object.
	ObjectTypeBlob ObjectType = "blob"
	// ObjectTypeTree represents a directory structure object.
	ObjectTypeTree ObjectType = "tree"
	// ObjectTypeCommit represents a commit snapshot object.
	ObjectTypeCommit ObjectType = "commit"
)

// IsValid checks if the ObjectType is one of the predefined valid types: blob, tree, or commit.
func (objectType ObjectType) IsValid() bool {
	switch objectType {
	case ObjectTypeBlob, ObjectTypeTree, ObjectTypeCommit:
		return true
	default:
		return false
	}
}

// String converts the ObjectType to its underlying string representation.
func (objectType ObjectType) String() string {
	return string(objectType)
}

// ParseObjectType converts a string to an ObjectType and verifies its validity, returning the type and a success flag.
func ParseObjectType(typeStr string) (ObjectType, bool) {
	objectType := ObjectType(typeStr)
	if objectType.IsValid() {
		return objectType, true
	}
	return "", false
}
