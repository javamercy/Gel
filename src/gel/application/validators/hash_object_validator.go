package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
	"Gel/src/gel/domain/objects"
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
		Must(isValidObjectType, "objectType must be one of Blob, Tree, Commit")

	fluentValidator.
		RuleFor("Paths", request.Paths).
		Must(isStringSliceNonEmpty, "path must be provided").
		Must(areAllInStringSliceNonEmpty, "paths must be non-empty")

	return fluentValidator.Validate()
}

func isValidObjectType(value any) bool {
	objectType, ok := value.(objects.ObjectType)
	return ok && objectType.IsValid()
}
