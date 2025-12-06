package logger

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var tp *sdktrace.TracerProvider

// InitTracer initializes OpenTelemetry tracer
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	// Get service name from environment if available
	if envServiceName := os.Getenv("OTEL_SERVICE_NAME"); envServiceName != "" {
		serviceName = envServiceName
	}

	// Get OTLP endpoint from environment, default to Jaeger
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "http://localhost:4318"
	}

	// Parse URL to extract host:port for WithEndpoint
	// WithEndpoint expects format like "localhost:4318" (without http://)
	var endpointHost string
	if parsedURL, err := url.Parse(otlpEndpoint); err == nil {
		endpointHost = parsedURL.Host
		if endpointHost == "" {
			// If parsing fails, try to extract host:port manually
			endpointHost = strings.TrimPrefix(strings.TrimPrefix(otlpEndpoint, "http://"), "https://")
		}
	} else {
		// Fallback: remove http:// or https:// prefix
		endpointHost = strings.TrimPrefix(strings.TrimPrefix(otlpEndpoint, "http://"), "https://")
	}

	// Create OTLP HTTP exporter for Jaeger
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpointHost),
		otlptracehttp.WithInsecure(), // Use insecure for local development
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	// Create resource with consistent schema URL
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Create tracer provider
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Printf("OpenTelemetry tracer initialized for service: %s", serviceName)
	return tp, nil
}

// Shutdown gracefully shuts down the tracer provider
func Shutdown() error {
	if tp != nil {
		return tp.Shutdown(context.Background())
	}
	return nil
}

// GetTracer returns a tracer instance
func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

