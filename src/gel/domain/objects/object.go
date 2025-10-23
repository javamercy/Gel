package objects

import (
	"Gel/src/gel/core/constant"
)

type IObject interface {
	Type() constant.ObjectType
	Size() int
	Data() []byte
}
type BaseObject struct {
	objectType constant.ObjectType
	data       []byte
}

func (baseObject *BaseObject) Type() constant.ObjectType {
	return baseObject.objectType
}

func (baseObject *BaseObject) Size() int {
	return len(baseObject.data)
}

func (baseObject *BaseObject) Data() []byte {
	return baseObject.data
}
