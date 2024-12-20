package validators

import (
	"github.com/go-playground/validator/v10"
)

func CustomErrorMessage(err validator.FieldError) string {
	fieldName := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fieldName + " is required"
	case "email":
		return fieldName + " must be a valid email"
	case "gte":
		return fieldName + " must be greater than or equal to " + param
	case "lte":
		return fieldName + " must be less than or equal to " + param
	default:
		return fieldName + " is invalid"
	}
}
