package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/validation"
	"regexp"
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
		Must(hashFormatMustBeValid, "Hash must be a valid hexadecimal string").
		Must(hashLengthMustBeValid, "Hash must be a valid SHA-256 hash")

	return fluentValidator.Validate()
}

func hashFormatMustBeValid(hash string) bool {
	matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", hash)
	return matched
}

func hashLengthMustBeValid(hash string) bool {
	return len(hash) == constant.SHA256HexLength
}
