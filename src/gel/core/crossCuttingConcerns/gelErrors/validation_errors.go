package gelErrors

type ValidationError struct {
	*GelError
	Field  string
	Detail string
}

func NewValidationError(field, detail, message string) *ValidationError {
	return &ValidationError{
		GelError: NewGelError(ExitCodeFatal, message),
		Field:    field,
		Detail:   detail,
	}
}
