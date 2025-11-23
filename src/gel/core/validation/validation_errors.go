package validation

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

func (error *ValidationError) Error() string {
	return error.Field + ": " + error.Message
}
