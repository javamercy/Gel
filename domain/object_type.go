package domain

import (
	"Gel/core/constant"
	"errors"
	"fmt"
)

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
		return ObjectTypeBlob, nil
	case constant.GelDirectoryModeStr:
		return ObjectTypeTree, nil
	default:
		return "", errors.New(fmt.Sprintf("Failed to get object type by mode: %v", mode))
	}
}
