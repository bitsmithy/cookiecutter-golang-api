package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	tracing "go.opentelemetry.io/otel/trace"
)

func setupTracer(ctx context.Context, r *resource.Resource) (*trace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(collectorURL)}
	if collectorURL == localCollectorURL {
		// We're in dev mode, running the collector with Docker
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)

	otel.SetTracerProvider(traceProvider)
	return traceProvider, nil
}

func Trace(ctx context.Context, name string) (context.Context, tracing.Span) {
	return otel.Tracer("").Start(ctx, name)
}

func CurrentSpan(ctx context.Context) tracing.Span {
	return tracing.SpanFromContext(ctx)
}

func WithSpanEvent(span tracing.Span, name string, myFunc func()) {
	span.AddEvent(name)
	defer span.AddEvent(fmt.Sprintf("%s done", name))

	myFunc()
}
