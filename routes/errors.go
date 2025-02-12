package routes

import (
	"coffee-shop-api/internal/common"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	Status  int
	Message interface{}
}

func FromError(ctx *fiber.Ctx, err error) error {
	var apiError APIError
	var svcError common.AError
	if errors.As(err, &svcError) {
		apiError.Message = ErrorResponse(svcError.AppError())
		apiError.Status = svcError.SvcError().Code
	}

	return ctx.Status(apiError.Status).JSON(apiError.Message)
}
