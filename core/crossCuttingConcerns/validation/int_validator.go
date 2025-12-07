package validation

type IntValidator struct {
	*FieldValidator
	value int
}

func NewIntValidator(fieldValidator *FieldValidator, value int) *IntValidator {
	return &IntValidator{
		FieldValidator: fieldValidator,
		value:          value,
	}
}
