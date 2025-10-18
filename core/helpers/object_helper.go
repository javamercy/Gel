package helpers

import (
	"Gel/core/constants"
	"strconv"
)

// PrepareObjectContent prepares the full object content in Git's format: <type> <size>\0<content>
func PrepareObjectContent(objectType constants.ObjectType, fileData []byte) []byte {
	header := string(objectType) + constants.Space + strconv.Itoa(len(fileData)) + constants.NullByte
	return append([]byte(header), fileData...)
}
