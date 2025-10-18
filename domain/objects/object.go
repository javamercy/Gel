package objects

import "Gel/core/constants"

type IObject interface {
	Type() constants.ObjectType
	Size() int
	Data() []byte
}
type BaseObject struct {
	objectType constants.ObjectType
	data       []byte
}

func (baseObject *BaseObject) Type() constants.ObjectType {
	return baseObject.objectType
}

func (baseObject *BaseObject) Size() int {
	return len(baseObject.data)
}

func (baseObject *BaseObject) Data() []byte {
	return baseObject.data
}
