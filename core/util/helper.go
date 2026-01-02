package util

import (
	"Gel/core/constant"
)

func FindNullByteIndex(data []byte) int {
	for i, b := range data {
		if b == constant.NullByte {
			return i
		}
	}
	return -1
}
