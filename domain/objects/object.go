package objects

import "Gel/core/constants"

type IObject interface {
	GetType() constants.ObjectType
	GetSize() int
	GetData() []byte
}
