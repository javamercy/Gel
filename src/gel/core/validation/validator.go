package validation

type IValidator interface {
	Validate(request any) *ValidationError
}
