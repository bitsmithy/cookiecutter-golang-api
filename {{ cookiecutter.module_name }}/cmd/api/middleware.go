package main

import (
  "context"
	"fmt"
	"net/http"

	"github.com/tomasen/realip"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"{{ cookiecutter.module_path }}/internal/log"
	"{{ cookiecutter.module_path }}/internal/response"
	"{{ cookiecutter.module_path }}/internal/telemetry"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorCounter, err := telemetry.Meter().Int64Counter(
			"error.counter",
			metric.WithDescription("Number of errors intercepted."),
			metric.WithUnit("{error}"),
		)
		if err != nil {
			app.serverError(w, r, fmt.Errorf("%w", err))
		}

    ctx := r.Context()
		span := telemetry.CurrentSpan(ctx)
    errorCounter.Add(r.Context(), 1)

		defer func(ctx context.Context) {
			err := recover()
			if err != nil {
				span.AddEvent("recoverPanic: error")
				app.serverError(w, r, fmt.Errorf("%s", err))
				errorCounter.Add(ctx, 1)
			}
		}(ctx)

		telemetry.WithSpanEvent(span, "recoverPanic: next", func() {
			next.ServeHTTP(w, r)
		})
	})
}

func (app *application) logAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCounter, err := telemetry.Meter().Int64Counter(
			"request.counter",
			metric.WithDescription("Number of API requests."),
			metric.WithUnit("{request}"),
		)
		if err != nil {
			app.serverError(w, r, fmt.Errorf("%w", err))
		}

		span := telemetry.CurrentSpan(r.Context())

		mw := response.NewMetricsResponseWriter(w)
		telemetry.WithSpanEvent(span, "logAccess: next", func() {
			next.ServeHTTP(mw, r)
		})

		attrs := []attribute.KeyValue{
			attribute.String("user.ip", realip.FromRequest(r)),
			attribute.String("request.method", r.Method),
			attribute.String("request.url", r.URL.String()),
			attribute.String("request.proto", r.Proto),
			attribute.Int("response.status", mw.StatusCode),
			attribute.Int("response.size", mw.BytesCount),
		}

		requestCounter.Add(r.Context(), 1, metric.WithAttributes(attrs...))
		log.
			WithSpanAttrs(span, attrs...).
			Info("request received")
	})
}

func (app *application) instrument(next http.Handler) http.Handler {
	instrument := func(w http.ResponseWriter, r *http.Request) {
		h := otelhttp.NewHandler(otelhttp.WithRouteTag(r.URL.Path, next), "server.http")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(instrument)
}
