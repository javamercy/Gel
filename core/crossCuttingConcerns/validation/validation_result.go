package validation

import "strings"

type ValidationResult struct {
	Errors []*ValidationError
}

func NewValidationResult(errors []*ValidationError) *ValidationResult {
	return &ValidationResult{
		Errors: errors,
	}
}

func (validationResult *ValidationResult) IsValid() bool {
	return len(validationResult.Errors) == 0
}

func (validationResult *ValidationResult) Error() string {
	errMessage := strings.Builder{}
	for _, err := range validationResult.Errors {
		errMessage.WriteString(err.Error())
		errMessage.WriteString("\n")
	}
	return strings.TrimSpace(errMessage.String())
}
