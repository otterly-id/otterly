package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidatorErrors(err error) []string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	errorMessages := []string{}

	for _, fieldErr := range validationErrors {
		var message string
		field := fieldErr.Field()
		param := fieldErr.Param()

		switch fieldErr.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long", field, param)
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters long", field, param)
		case "gte":
			message = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
		case "lte":
			message = fmt.Sprintf("%s must be less than or equal to %s characters", field, param)
		case "oneof":
			message = fmt.Sprintf("%s must be one of the following: %s", field, strings.ReplaceAll(param, " ", ", "))
		case "uuid":
			message = fmt.Sprintf("%s must be a valid UUID", field)
		case "alpha_space":
			message = fmt.Sprintf("%s must contain only letters and spaces", field)
		case "phone":
			message = fmt.Sprintf("%s must be a valid phone number (e.g., +1234567890)", field)
		case "password_strength":
			message = fmt.Sprintf("%s must be at least 8 characters long and contain at least 1 uppercase letter, 1 lowercase letter, and 1 number", field)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		errorMessages = append(errorMessages, message)
	}

	return errorMessages
}
