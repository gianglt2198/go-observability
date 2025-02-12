package middlewares

import (
	"coffee-shop-api/monitoring"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig holds configuration for the tracing middleware
type TracingConfig struct {
	// ServiceName is the name of the service being traced
	ServiceName string
	// ServiceVersion is the version of the service
	ServiceVersion string
	// Skip defines a function to skip the middleware
	Skip func(c *fiber.Ctx) bool
	// Whitelist path to skip the middleware
	Whitelist map[string]bool
}

// DefaultTracingConfig returns the default configuration
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		ServiceName:    "coffee-shop-api",
		ServiceVersion: "1.0.0",
		Skip:           nil,
		Whitelist: map[string]bool{
			"/health":  true,
			"/metrics": true,
			"/swagger": true,
		},
	}
}

func (m *TracingConfig) apply(cfg *TracingConfig) {
	if m.ServiceName != "" {
		cfg.ServiceName = m.ServiceName
	}
	if m.ServiceVersion != "" {
		cfg.ServiceVersion = m.ServiceVersion
	}
	if m.Skip != nil {
		cfg.Skip = m.Skip
	}
	if m.Whitelist != nil {
		cfg.Whitelist = m.Whitelist
	}
}

// TracingMiddleware returns a middleware that traces HTTP requests using OpenTelemetry
func TracingMiddleware(handlerName, operationName string, configs ...TracingConfig) fiber.Handler {
	// Use default config if none provided
	cfg := DefaultTracingConfig()
	if len(configs) > 0 {
		for _, c := range configs {
			c.apply(&cfg)
		}
	}

	return func(c *fiber.Ctx) error {

		// Check if middleware should be skipped
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}

		// Start span
		ctx := c.Context()
		tracer := monitoring.Tracer(c.OriginalURL())
		spanName := fmt.Sprintf("%s %s", c.Method(), c.OriginalURL())

		tracingCtx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.url", c.OriginalURL()),
				attribute.String("http.route", c.Route().Path),
				attribute.String("http.host", c.Hostname()),
				attribute.String("http.request_id", c.Locals("requestId").(string)),
				attribute.String("service.name", cfg.ServiceName),
				attribute.String("service.version", cfg.ServiceVersion),
				attribute.String("http.user_agent", string(c.Request().Header.UserAgent())),
			),
		)
		defer span.End()

		// Store span context in fiber context
		c.Locals("span", span)
		c.Locals("spanCtx", tracingCtx)

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Add response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
			attribute.String("http.response_content_type", c.Get(fiber.HeaderContentType)),
			attribute.Float64("http.request_duration_ms", float64(time.Since(start).Milliseconds())),
		)

		// Handle errors
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err, trace.WithAttributes(
				attribute.String("error.type", fmt.Sprintf("%T", err)),
				attribute.String("error.message", err.Error()),
			))
			return err
		}

		// Set status based on response code
		if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			span.SetStatus(codes.Ok, "requset successed")
		} else {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", c.Response().StatusCode()))
		}

		return err
	}
}

// GetSpanFromContext retrieves the current span from the fiber context
func GetSpanFromContext(c context.Context, spanName string) (context.Context, trace.Span) {
	if ctx, ok := c.Value("spanCtx").(context.Context); ok {
		c = ctx
	}
	childCtx, childSpan := monitoring.Tracer(spanName).Start(c, spanName)
	return childCtx, childSpan
}
