package middlewares

import (
	"coffee-shop-api/monitoring"
	"context"
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MetricConfig holds configuration for the tracing middleware
type MetricConfig struct {
	// ServiceName is the name of the service being traced
	ServiceName string
	// ServiceVersion is the version of the service
	ServiceVersion string
	// Skip defines a function to skip the middleware
	Skip func(c *fiber.Ctx) bool
	// Whitelist path to skip the middleware
	Whitelist map[string]bool
	// Metrics list to track requests
	// *** counter ***
	httpRequestsCounter       metric.Int64Counter
	httpActiveRequestsCounter metric.Int64UpDownCounter

	// *** gauge ***
	memoryUsageObservableGuage metric.Int64ObservableGauge

	// *** histogram ***
	httpRequestDurationHistogram metric.Int64Histogram
}

// DefaultMetricConfig returns the default configuration
func DefaultMetricConfig() MetricConfig {
	metricConfig := MetricConfig{
		ServiceName:    "coffee-shop-api",
		ServiceVersion: "1.0.0",
		Skip:           nil,
		Whitelist: map[string]bool{
			"/health":  true,
			"/metrics": true,
			"/swagger": true,
		},
	}

	metricConfig.newCounters()
	metricConfig.newHistograms()
	metricConfig.newGauges()

	return metricConfig
}

func (m MetricConfig) apply(cfg *MetricConfig) {
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

func MetricMiddleware(configs ...MetricConfig) fiber.Handler {
	// Use default config if none provided
	cfg := DefaultMetricConfig()
	if len(configs) > 0 {
		for _, c := range configs {
			c.apply(&cfg)
		}
	}

	return func(c *fiber.Ctx) error {

		if cfg.Skip != nil && cfg.Skip(c) || cfg.Whitelist[c.OriginalURL()] {
			return c.Next()
		}

		start := time.Now().UTC()

		ctx := c.Context()

		metricAttributes := attribute.NewSet(
			attribute.String("http.method", c.Method()),
			attribute.String("http.url", c.OriginalURL()),
			attribute.String("http.route", c.Route().Path),
			attribute.String("http.host", c.Hostname()),
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
		)

		cfg.httpRequestsCounter.Add(
			ctx, 1, metric.WithAttributeSet(metricAttributes),
		)

		cfg.httpActiveRequestsCounter.Add(ctx, 1, metric.WithAttributeSet(metricAttributes))

		// Process
		err := c.Next()

		// Recording Metric
		cfg.httpActiveRequestsCounter.Add(ctx, -1, metric.WithAttributeSet(metricAttributes))
		cfg.httpRequestDurationHistogram.Record(ctx,
			time.Since(start).Milliseconds(),
			metric.WithAttributeSet(metricAttributes))

		if err != nil {
			return err
		}

		return nil
	}
}

func (c *MetricConfig) newCounters() {
	m := monitoring.Meter(c.ServiceName)

	var err error

	c.httpRequestsCounter, err = m.Int64Counter(
		"http_service_requests_total",
		metric.WithDescription("Total number of HTTP requests received."),
		metric.WithUnit("{requests}"),
	)
	if err != nil {
		log.Fatalf("creating meter http request counter failed: %v", err)
	}

	c.httpActiveRequestsCounter, err = m.Int64UpDownCounter(
		"http_service_active_requests_total",
		metric.WithDescription("Number of in-flight requests."),
		metric.WithUnit("{requests}"),
	)
	if err != nil {
		log.Fatalf("creating meter http active request counter failed: %v", err)
	}

}

func (c *MetricConfig) newGauges() {
	m := monitoring.Meter(c.ServiceName)

	var err error

	c.memoryUsageObservableGuage, err = m.Int64ObservableGauge(
		"system_memory_heap",
		metric.WithDescription(
			"Memory usage of the allocated heap objects.",
		),
		metric.WithUnit("By"),
		metric.WithInt64Callback(
			func(ctx context.Context, o metric.Int64Observer) error {
				var memStats runtime.MemStats

				runtime.ReadMemStats(&memStats)

				currentMemoryUsage := memStats.HeapAlloc

				o.Observe(int64(currentMemoryUsage))
				return nil
			},
		),
	)
	if err != nil {
		log.Fatalf("creating meter memory usage gauge failed: %v", err)
	}
}

func (c *MetricConfig) newHistograms() {
	m := monitoring.Meter(c.ServiceName)

	var err error

	c.httpRequestDurationHistogram, err = m.Int64Histogram(
		"http_service_request_duration_milliseconds",
		metric.WithDescription("The duration of an HTTP request."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalf("creating meter http request duration failed: %v", err)
	}

}
