package validators

import (
	"Gel/src/gel/application/dto"
	"Gel/src/gel/core/validation"
	"path/filepath"
	"strings"
)

type InitValidator struct {
}

func NewInitValidator() *InitValidator {
	return &InitValidator{}
}

func (initValidator *InitValidator) Validate(data any) *validation.ValidationError {
	request, ok := data.(dto.InitRequest)
	if !ok {
		return validation.NewValidationError("request", "invalid request type")
	}

	if pathMustNotBeEmpty(request.Path) {
		return validation.NewValidationError("path", "path must not be empty")
	}

	if isDangerousPath(request.Path) {
		return validation.NewValidationError("path", "path is too dangerous to initialize a repository")
	}
	return nil
}

func isDangerousPath(path string) bool {
	dangerousPaths := []string{"/", "/etc", "/usr", "/bin", "/sbin", "/var", "/tmp", "/boot"}
	absPath, _ := filepath.Abs(path)
	for _, path := range dangerousPaths {
		if absPath == path || strings.HasPrefix(absPath, path+"/") {
			return true
		}
	}
	return false
}

func pathMustNotBeEmpty(path string) bool {
	return strings.TrimSpace(path) == ""
}
