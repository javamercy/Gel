package validation

type ArrayValidator struct {
	*FieldValidator
	value []any
}

func NewArrayValidator(fieldValidator *FieldValidator, value []any) *ArrayValidator {
	return &ArrayValidator{
		FieldValidator: fieldValidator,
		value:          value,
	}
}

func (arrayValidator *ArrayValidator) NotEmpty() *ArrayValidator {
	if len(arrayValidator.value) == 0 {
		validationError := NewValidationError(arrayValidator.fieldName, "must not be empty")
		arrayValidator.parent.AddError(validationError)
	}
	return arrayValidator
}

func (arrayValidator *ArrayValidator) Must(predicate func([]any) bool, message string) *ArrayValidator {
	if arrayValidator.stop() {
		return arrayValidator
	}

	if !predicate(arrayValidator.value) {
		validationError := NewValidationError(
			arrayValidator.fieldName,
			message,
		)
		arrayValidator.parent.AddError(validationError)
	}
	return arrayValidator
}

func (arrayValidator *ArrayValidator) WithMessage(message string) *ArrayValidator {
	lastIndex := len(arrayValidator.parent.Errors) - 1
	if lastIndex >= 0 {
		arrayValidator.parent.Errors[lastIndex].Message = message
	}
	return arrayValidator
}
