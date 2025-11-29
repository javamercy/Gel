package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
)

type UpdateIndexValidator struct {
}

func NewUpdateIndexValidator() *UpdateIndexValidator {
	return &UpdateIndexValidator{}
}

func (updateIndexValidator *UpdateIndexValidator) Validate(request *dto.UpdateIndexRequest) *validation.ValidationResult {
	fluentValidator := validation.NewFluentValidator(true)

	fluentValidator.RuleFor("Paths", request.Paths).
		Must(isStringSliceNonEmpty, "paths must contain at least one path").
		Must(areAllInStringSliceNonEmpty, "all paths must be non-empty strings").
		Must(isOneOfOptions, "at least one of options must be used").
		Must(isOnlyOneOption, "only one of options can be used at a time")

	return fluentValidator.Validate()
}

func isOneOfOptions(value any) bool {
	request, ok := value.(*dto.UpdateIndexRequest)
	if !ok {
		return false
	}
	return atLeastOne(request.Add, request.Remove)
}

func isOnlyOneOption(value any) bool {
	request, ok := value.(*dto.UpdateIndexRequest)
	if !ok {
		return false
	}
	return exactlyOne(request.Add, request.Remove)
}
