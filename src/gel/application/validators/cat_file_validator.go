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
		Matches(regexSHA256).
		WithMessage("hash must be a valid SHA-256 hash")

	fluentValidator.
		RuleFor("Options", request).
		Must(atLeastOneCatFileOption, "at least one of options must be used").
		Must(isOnlyOneCatFileOption, "only one of options can be used at a time")

	return fluentValidator.Validate()
}

func atLeastOneCatFileOption(value any) bool {
	request, ok := value.(*dto.CatFileRequest)
	if !ok {
		return false
	}
	return atLeastOne(request.ShowType, request.ShowSize, request.Pretty, request.CheckOnly)
}

func isOnlyOneCatFileOption(value any) bool {
	request, ok := value.(*dto.CatFileRequest)
	if !ok {
		return false
	}
	return exactlyOne(request.ShowType, request.ShowSize, request.Pretty, request.CheckOnly)
}
