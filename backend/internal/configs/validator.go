package configs

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v.RegisterValidation("alpha_space", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(fl.Field().String())
	})

	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if field == "" {
			return true
		}
		return regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(field)
	})

	v.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		return len(password) >= 8 && hasUpper && hasLower && hasNumber
	})

	return v
}
