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
		RuleFor("ObjectType", request.ObjectType).
		Must(isValidObjectType, "ObjectType must be one of Blob, Tree, Commit")

	fluentValidator.
		RuleFor("Paths", request.Paths).
		Must(isStringSliceNonEmpty, "Paths must contain at least one path").
		Must(areAllInStringSliceNonEmpty, "All paths must be non-empty strings")

	return fluentValidator.Validate()
}

func isValidObjectType(value any) bool {
	objectType, ok := value.(constant.ObjectType)
	return ok && (objectType == constant.GelBlobObjectType || objectType == constant.GelTreeObjectType || objectType == constant.GelCommitObjectType)
}
