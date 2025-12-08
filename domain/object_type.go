package domain

import (
	"Gel/core/constant"
	"errors"
	"fmt"
)

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

func GetObjectTypeByMode(mode string) (ObjectType, error) {
	switch mode {
	case constant.GelRegularFileModeStr, constant.GelExecutableFileModeStr:
		return GelBlobObjectType, nil
	case constant.GelDirectoryModeStr:
		return GelTreeObjectType, nil
	default:
		return "", errors.New(fmt.Sprintf("Failed to get object type by mode: %v", mode))
	}
}
