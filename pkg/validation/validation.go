package validation

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Register Tag Json as field name not from struct name for validation errors
func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.Split(fld.Tag.Get("json"), ",")[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
	}
}

func FormatValidationErrors(ve validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, e := range ve {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errs[field] = field + " is required"
		case "url":
			errs[field] = field + " must be a valid URL"
		case "min":
			errs[field] = field + " must be at least " + e.Param()
		case "max":
			errs[field] = field + " must be at most " + e.Param()
		case "email":
			errs[field] = field + " must be a valid email address"
		case "alphanum":
			errs[field] = field + " must only contains letters and numbers"
		default:
			errs[field] = field + " is invalid"
		}
	}
	return errs
}
