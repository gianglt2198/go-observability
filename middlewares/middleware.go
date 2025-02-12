package middlewares

import (
	"coffee-shop-api/routes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func AllPayloadValidator[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var data T
		var err error

		_ = c.BodyParser(&data)

		_ = c.ParamsParser(&data)

		_ = c.QueryParser(&data)

		err = routes.ValidateStruct(&data)

		if err != nil {
			return err
		}

		c.Locals(routes.KEY_REQ_ALL_PARAMS, data)

		return c.Next()
	}
}

func ParseValidationError(err error) string {
	if err == nil {
		return ""
	}

	errors := make(map[string]interface{})

	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}

	toJson, _ := json.Marshal(errors)

	return string(toJson)
}
