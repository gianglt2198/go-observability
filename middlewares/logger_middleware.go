package middlewares

import (
	"coffee-shop-api/monitoring"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	REQUEST_ID = "requestId"
)

func RequestIDMiddleware(c *fiber.Ctx) error {
	requestID := c.Locals(REQUEST_ID)
	if s, ok := requestID.(string); !ok && s == "" {
		requestId := uuid.New().String()
		c.Locals(REQUEST_ID, requestId)
	}

	return c.Next()
}

func LoggerMiddleware(logger *monitoring.AppLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		logFields := []zapcore.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("request_id", c.Locals("requestId").(string)),
		}
		if err := c.Next(); err != nil {
			logger.GetLogger().With(logFields...).Error("Request failed", zap.Error(err))
			return err
		}

		logger.GetLogger().With(logFields...).Info("Request succeeded")
		return nil
	}
}
