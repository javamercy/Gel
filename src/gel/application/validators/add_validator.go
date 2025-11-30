package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
)

type AddValidator struct {
}

func NewAddValidator() *AddValidator {
	return &AddValidator{}
}

func (addValidator *AddValidator) Validate(request *dto.AddRequest) *validation.ValidationResult {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.RuleFor("Pathspecs", request.Pathspecs).
		Must(isStringSliceNonEmpty, "at least one pathspec must be provided").
		Must(areAllInStringSliceNonEmpty, "all pathspecs must be non-empty strings")

	return fluentValidator.Validate()
}
