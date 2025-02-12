package monitoring

import (
	"coffee-shop-api/config"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	metric2 "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Tracer(name string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name)
}

func Meter(name string) metric2.Meter {
	return otel.GetMeterProvider().Meter(name)
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, cfg *config.Config) (shutdown func(context.Context) error, err error) {

	log.Println("Initializing OpenTelemetry SDK")

	// Verbose error handling
	handleErr := func(inErr error) {
		log.Printf("OpenTelemetry setup error: %v", inErr)
		err = errors.Join(inErr, shutdown(ctx))
	}

	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ctx, cfg)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, cfg)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// // Set up logger provider.
	// loggerProvider, err := newLoggerProvider(ctx, cfg)
	// if err != nil {
	// 	handleErr(err)
	// 	return
	// }
	// shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	// global.SetLoggerProvider(loggerProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context, cfg *config.Config) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Tracing.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create a resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.App.Name),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating trace resource: %w", err)
	}

	// Create trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return traceProvider, nil
}

func newMeterProvider(ctx context.Context, cfg *config.Config) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.Tracing.Endpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithCompressor("gzip"),
	)
	// metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	// Create a resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.App.Name),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating metric resource: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(5*time.Second))),
		metric.WithResource(res),
	)
	return meterProvider, nil
}

// func newLoggerProvider(ctx context.Context, cfg *config.Config) (*log.LoggerProvider, error) {
// 	logExporter, err := stdoutlog.New()
// 	if err != nil {
// 		return nil, err
// 	}

// 	loggerProvider := log.NewLoggerProvider(
// 		log.WithProcessor(log.NewBatchProcessor(logExporter)),
// 	)
// 	return loggerProvider, nil
// }
