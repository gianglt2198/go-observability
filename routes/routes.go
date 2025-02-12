package routes

import (
	"coffee-shop-api/internal/common"
	"context"

	"github.com/gofiber/fiber/v2"
)

const (
	KEY_REQ_ALL_PARAMS = "parsed_all_params"
)

type Handler[T any, R any] func(context.Context, T) (R, error)

func Usecase[T any, R any](f Handler[T, R], successStatus int, ms ...fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Apply middlewares in reverse order
		for i := len(ms) - 1; i >= 0; i-- {
			middleware := ms[i]
			if err := middleware(ctx); err != nil {
				return FromError(ctx, common.NewError(fiber.ErrBadRequest, err))
			}
		}

		// Get parsed data from context
		data, ok := ctx.Locals(KEY_REQ_ALL_PARAMS).(T)
		if !ok {
			return FromError(ctx, common.NewError(fiber.ErrBadRequest, common.ErrBadRequest))
		}

		// Execute handler
		result, err := f(ctx.Context(), data)
		if err != nil {
			return FromError(ctx, err)
		}

		return ctx.Status(successStatus).JSON(SuccessResponse(result))
	}
}
