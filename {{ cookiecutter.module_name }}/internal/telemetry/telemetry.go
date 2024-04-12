package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"

	"{{ cookiecutter.module_path }}/internal/env"
	"{{ cookiecutter.module_path }}/internal/log"
	"{{ cookiecutter.module_path }}/internal/version"

	semconv "go.opentelemetry.io/otel/semconv/v1.23.1"
)

const localCollectorURL = "localhost:4317"

var collectorURL = env.GetString("OTEL_COLLECTOR_URL", localCollectorURL)

// Setup bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func Setup(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	log.With("endpoint", collectorURL).Info("setting up telemetry")

	var shutdownFuncs []func(context.Context) error
	var err error

	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		log.Info("shutting down telemetry")
		return nil
	})

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	setupPropagator()

	// Set up a default resource
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersion(version.Get()),
		semconv.DeploymentEnvironment(env.GetString("TELEMETRY_ENV", "local")),
	)

	// Set up trace provider.
	tracerProvider, err := setupTracer(ctx, resource)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// Set up meter provider.
	meterProvider, err := setupMeter(ctx, resource)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return shutdown, err
}
