package validation

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func (validationError *ValidationError) Error() string {
	return fmt.Sprintf("Field '%v': %v", validationError.Field, validationError.Message)
}
