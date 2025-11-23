package dto

import "Gel/src/gel/core/constant"

type InitRequest struct {
	Path string
}

func NewInitRequest(path string) *InitRequest {
	return &InitRequest{
		path,
	}
}

type HashObjectRequest struct {
	Paths      []string
	ObjectType constant.ObjectType
	Write      bool
}

func NewHashObjectRequest(paths []string, objectType constant.ObjectType, write bool) *HashObjectRequest {
	return &HashObjectRequest{
		Paths:      paths,
		ObjectType: objectType,
		Write:      write,
	}
}
