package validation

type IValidator interface {
	Validate(data any) *ValidationError
}
