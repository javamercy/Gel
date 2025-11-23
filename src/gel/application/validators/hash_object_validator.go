package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/validation"
)

type HashObjectValidator struct {
}

func NewHashObjectValidator() *HashObjectValidator {
	return &HashObjectValidator{}
}

func (hashObjectValidator *HashObjectValidator) Validate(request any) *validation.ValidationError {

	hashObjectRequest, ok := request.(*dto.HashObjectRequest)
	if !ok {
		return validation.NewValidationError("request", "Invalid request type")
	}

	if !pathsMustNotBeEmpty(hashObjectRequest.Paths) {
		return validation.NewValidationError("paths", "Paths must not be empty")
	}

	if !objectTypeMustBeValid(hashObjectRequest.ObjectType) {
		return validation.NewValidationError("objectType", "Invalid object type")
	}

	return nil
}

func pathsMustNotBeEmpty(paths []string) bool {
	return len(paths) > 0
}

func objectTypeMustBeValid(objectType constant.ObjectType) bool {
	switch objectType {
	case constant.GelBlobObjectType, constant.GelTreeObjectType, constant.GelCommitObjectType:
		return true
	default:
		return false
	}
}
