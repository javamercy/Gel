package objects

import "Gel/src/gel/core/constant"

type ObjectType string

const (
	GelBlobObjectType   ObjectType = "blob"
	GelTreeObjectType   ObjectType = "tree"
	GelCommitObjectType ObjectType = "commit"
)

func (objectType ObjectType) IsValid() bool {
	switch objectType {
	case GelBlobObjectType, GelTreeObjectType, GelCommitObjectType:
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

func GetObjectTypeByMode(mode string) (string, bool) {
	switch mode {
	case constant.GelRegularFileMode, constant.GelExecFileMode:
		return GelBlobObjectType.String(), true
	case constant.GelDirMode:
		return GelTreeObjectType.String(), true
	default:
		return "", false
	}
}
