package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
)

func setupMeter(ctx context.Context, r *resource.Resource) (*sdkMetric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(collectorURL)}
	if collectorURL == localCollectorURL {
		// We're in dev mode, running the collector with Docker
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(exporter)),
		sdkMetric.WithResource(r),
	)

	otel.SetMeterProvider(meterProvider)
	return meterProvider, nil
}

func Meter() metric.Meter {
	return otel.Meter("")
}
