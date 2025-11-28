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

func (hashObjectValidator *HashObjectValidator) Validate(request *dto.HashObjectRequest) *validation.ValidationResult {

	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.
		RuleFor("Paths", request.Paths).
		Array().
		NotEmpty()

	return fluentValidator.Validate()
}

func objectTypeMustBeValid(objectType constant.ObjectType) bool {
	switch objectType {
	case constant.GelBlobObjectType, constant.GelTreeObjectType, constant.GelCommitObjectType:
		return true
	default:
		return false
	}
}
