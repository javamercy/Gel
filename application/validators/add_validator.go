package validators

import (
	"Gel/application/dto"
	"Gel/core/crossCuttingConcerns/gelErrors"
	"Gel/core/crossCuttingConcerns/validation"
)

type AddValidator struct {
}

func NewAddValidator() *AddValidator {
	return &AddValidator{}
}

func (addValidator *AddValidator) Validate(request *dto.AddRequest) *gelErrors.GelError {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.RuleFor("Pathspecs", request.Pathspecs).
		Must(isStringSliceNonEmpty, "at least one pathspec must be provided").
		Must(areAllInStringSliceNonEmpty, "all pathspecs must be non-empty strings")

	validationResult := fluentValidator.Validate()
	if !validationResult.IsValid() {
		return gelErrors.NewGelError(gelErrors.ExitCodeUsage,
			validationResult.Error())
	}

	return nil
}
