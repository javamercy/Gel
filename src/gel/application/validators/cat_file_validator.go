package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/constant"
	"Gel/src/gel/core/validation"
	"regexp"
	"strings"
)

type CatFileValidator struct {
}

func NewCatFileValidator() *CatFileValidator {
	return &CatFileValidator{}
}

func (catFileValidator *CatFileValidator) Validate(data any) *validation.ValidationError {
	request, ok := data.(*dto.CatFileRequest)
	if !ok {
		return validation.NewValidationError("request", "Invalid request type")
	}

	if hashMustNotBeEmpty(request.Hash) {
		return validation.NewValidationError("hash", "Hash must not be empty")
	}

	if !hashFormatMustBeValid(request.Hash) {
		return validation.NewValidationError("hash", "Hash must contain only hexadecimal characters")
	}

	if !hashLengthMustBeValid(request.Hash) {
		return validation.NewValidationError("hash", "Hash must be exactly 64 characters (SHA-256)")
	}

	return nil
}

func hashMustNotBeEmpty(hash string) bool {
	return strings.TrimSpace(hash) == ""
}

func hashFormatMustBeValid(hash string) bool {
	matched, _ := regexp.MatchString("^[0-9a-fA-F]+$", hash)
	return matched
}

func hashLengthMustBeValid(hash string) bool {
	return len(hash) == constant.SHA256HexLength
}
