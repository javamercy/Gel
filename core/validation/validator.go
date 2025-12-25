package validation

import (
	"Gel/core/constant"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		_ = validate.RegisterValidation("sha256hex", validateSha256Hex)
		_ = validate.RegisterValidation("relativepath", validateRelativePath)
		_ = validate.RegisterValidation("timezone", validateTimezone)
	})
	return validate
}

func validateSha256Hex(fieldLevel validator.FieldLevel) bool {
	hash := fieldLevel.Field().String()
	matched, err := regexp.MatchString(`^[0-9a-f]{64}$`, hash)
	return err == nil && matched
}

func validateRelativePath(fieldLevel validator.FieldLevel) bool {
	path := fieldLevel.Field().String()
	if path == "" {
		return false
	}
	if path[0] == constant.SlashByte || strings.Contains(path, constant.PreviousDirectoryStr) {
		return false
	}
	return true
}
func validateTimezone(fieldLevel validator.FieldLevel) bool {
	timezone := fieldLevel.Field().String()
	matched, err := regexp.MatchString(`^[+-]\d{4}$`, timezone)
	return err == nil && matched
}
