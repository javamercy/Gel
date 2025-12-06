package validation

type FluentValidator struct {
	Errors        []*ValidationError
	StopOnFailure bool
}

func NewFluentValidator(stopOnFailure bool) *FluentValidator {
	return &FluentValidator{
		Errors:        []*ValidationError{},
		StopOnFailure: stopOnFailure,
	}
}

func (fluentValidator *FluentValidator) RuleFor(fieldName string, value any) *FieldValidator {
	return NewFieldValidator(fieldName, value, fluentValidator)
}

func (fluentValidator *FluentValidator) Validate() *ValidationResult {
	return NewValidationResult(fluentValidator.Errors)
}
func (fluentValidator *FluentValidator) HasErrors() bool {
	return len(fluentValidator.Errors) > 0
}

func (fluentValidator *FluentValidator) AddError(validationError *ValidationError) {
	fluentValidator.Errors = append(fluentValidator.Errors, validationError)
}
