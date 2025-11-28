package validation

import (
	"regexp"
	"strings"
)

type StringValidator struct {
	*FieldValidator
	value string
}

func NewStringValidator(fieldValidator *FieldValidator, value string) *StringValidator {
	return &StringValidator{
		FieldValidator: fieldValidator,
		value:          value,
	}
}

func (stringValidator *StringValidator) NotEmpty() *StringValidator {
	if stringValidator.stop() {
		return stringValidator
	}

	if strings.TrimSpace(stringValidator.value) == "" {
		ValidationError := NewValidationError(
			stringValidator.fieldName,
			"must not be empty",
		)
		stringValidator.parent.AddError(ValidationError)
	}

	return stringValidator
}

func (stringValidator *StringValidator) Matches(regexp regexp.Regexp) *StringValidator {
	if stringValidator.stop() {
		return stringValidator
	}
	if !regexp.MatchString(stringValidator.value) {
		ValidationError := NewValidationError(
			stringValidator.fieldName,
			"has invalid format",
		)
		stringValidator.parent.AddError(ValidationError)
	}
	return stringValidator
}

func (stringValidator *StringValidator) WithMessage(message string) *StringValidator {

	lastIndex := len(stringValidator.parent.Errors) - 1
	if lastIndex >= 0 {
		stringValidator.parent.Errors[lastIndex].Message = message
	}
	return stringValidator
}

func (stringValidator *StringValidator) Must(predicate func(string) bool, message string) *StringValidator {
	if stringValidator.stop() {
		return stringValidator
	}

	if !predicate(stringValidator.value) {
		ValidationError := NewValidationError(
			stringValidator.fieldName,
			message,
		)
		stringValidator.parent.AddError(ValidationError)
	}
	return stringValidator
}
