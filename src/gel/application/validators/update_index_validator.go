package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
)

type UpdateIndexValidator struct {
}

func NewUpdateIndexValidator() *UpdateIndexValidator {
	return &UpdateIndexValidator{}
}

func (updateIndexValidator *UpdateIndexValidator) Validate(request *dto.UpdateIndexRequest) *validation.ValidationResult {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.RuleFor("Paths", request.Paths).
		Must(isStringSliceNonEmpty, "Paths must contain at least one path").
		Must(areAllInStringSliceNonEmpty, "All paths must be non-empty strings")

	return fluentValidator.Validate()
}
