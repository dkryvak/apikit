package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var v = newValidator()

func Struct(s any) error {
	if err := v.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return format(ve)
		}
		return err
	}
	return nil
}

func newValidator() *validator.Validate {
	val := validator.New(validator.WithRequiredStructEnabled())

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Tag.Get("json")
		if tag == "" || tag == "-" {
			return ""
		}
		name := strings.Split(tag, ",")[0]
		if name == "" || name == "-" {
			return ""
		}
		return name
	})

	return val
}

func format(ve validator.ValidationErrors) error {
	var parts []string
	for _, e := range ve {
		field := e.Field()
		if field == "" {
			field = e.StructField()
		}

		switch e.Tag() {
		case "required":
			parts = append(parts, fmt.Sprintf("%s is required", field))
		case "oneof":
			parts = append(parts, fmt.Sprintf("%s must be one of [%s]", field, e.Param()))
		case "min":
			parts = append(parts, fmt.Sprintf("%s must be at least %s", field, e.Param()))
		case "max":
			parts = append(parts, fmt.Sprintf("%s must be at most %s", field, e.Param()))
		case "url":
			parts = append(parts, fmt.Sprintf("%s must be a valid URL", field))
		default:
			parts = append(parts, fmt.Sprintf("%s failed on '%s'", field, e.Tag()))
		}
	}
	return fmt.Errorf("validation failed: %s", strings.Join(parts, "; "))
}
