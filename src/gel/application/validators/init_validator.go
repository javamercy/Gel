package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/crossCuttingConcerns/gelErrors"
	"Gel/src/gel/core/validation"
)

type InitValidator struct {
}

func NewInitValidator() *InitValidator {
	return &InitValidator{}
}

func (initValidator *InitValidator) Validate(request *dto.InitRequest) *gelErrors.GelError {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.
		RuleFor("Path", request.Path).
		String().
		NotEmpty()

	validationResult := fluentValidator.Validate()
	if !validationResult.IsValid() {
		return gelErrors.NewGelError(gelErrors.ExitCodeUsage,
			validationResult.Error())
	}

	return nil
}
