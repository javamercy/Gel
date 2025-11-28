package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
)

type InitValidator struct {
}

func NewInitValidator() *InitValidator {
	return &InitValidator{}
}

func (initValidator *InitValidator) Validate(request *dto.InitRequest) *validation.ValidationResult {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.
		RuleFor("Path", request.Path).
		String().
		NotEmpty()

	return fluentValidator.Validate()
}
