package routes

import "github.com/go-playground/validator/v10"

func ValidateStruct(data interface{}) error {
	v := validator.New(validator.WithRequiredStructEnabled())

	return v.Struct(data)
}
