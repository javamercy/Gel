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

	fluentValidator.
		RuleFor("Paths", request.Paths).
		Array().
		NotEmpty()

	return fluentValidator.Validate()
}
