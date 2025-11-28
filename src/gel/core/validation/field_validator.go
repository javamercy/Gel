package validation

type FieldValidator struct {
	fieldName string
	value     any
	parent    *FluentValidator
	valid     bool
}

func NewFieldValidator(fieldName string, value any, parent *FluentValidator) *FieldValidator {
	return &FieldValidator{
		fieldName: fieldName,
		value:     value,
		parent:    parent,
		valid:     true,
	}
}

func (fieldValidator *FieldValidator) String() *StringValidator {
	stringValue, ok := fieldValidator.value.(string)
	fieldValidator.valid = ok

	if !ok {
		fieldValidator.parent.Errors = append(fieldValidator.parent.Errors, NewValidationError(
			fieldValidator.fieldName,
			"must be a string",
		))

		return NewStringValidator(fieldValidator, "")
	}

	return NewStringValidator(fieldValidator, stringValue)
}

func (fieldValidator *FieldValidator) Int() *IntValidator {
	intValue, ok := fieldValidator.value.(int)
	fieldValidator.valid = ok

	if !ok {
		fieldValidator.parent.Errors = append(fieldValidator.parent.Errors, NewValidationError(
			fieldValidator.fieldName,
			"must be an integer",
		))

		return NewIntValidator(fieldValidator, 0)
	}

	return NewIntValidator(fieldValidator, intValue)
}

func (fieldValidator *FieldValidator) Must(predicate func(any) bool, message string) *FieldValidator {
	if fieldValidator.stop() {
		return fieldValidator
	}

	if !predicate(fieldValidator.value) {
		validationError := NewValidationError(
			fieldValidator.fieldName,
			message,
		)
		fieldValidator.parent.AddError(validationError)
	}
	return fieldValidator
}

func (fieldValidator *FieldValidator) stop() bool {
	return !fieldValidator.valid || (fieldValidator.parent.StopOnFailure && fieldValidator.parent.HasErrors())
}
