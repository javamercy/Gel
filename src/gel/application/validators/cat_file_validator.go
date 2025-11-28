package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
)

type CatFileValidator struct {
}

func NewCatFileValidator() *CatFileValidator {
	return &CatFileValidator{}
}

func (catFileValidator *CatFileValidator) Validate(request *dto.CatFileRequest) *validation.ValidationResult {
	fluentValidator := validation.NewFluentValidator(false)

	fluentValidator.
		RuleFor("Hash", request.Hash).
		String().
		NotEmpty().
		Matches(RegexSHA256)

	return fluentValidator.Validate()
}
