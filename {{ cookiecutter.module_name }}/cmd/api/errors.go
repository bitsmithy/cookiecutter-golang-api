package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"{{ cookiecutter.module_path }}/internal/log"
	"{{ cookiecutter.module_path }}/internal/response"
	"{{ cookiecutter.module_path }}/internal/telemetry"
	"{{ cookiecutter.module_path }}/internal/validator"
)

func (app *application) reportServerError(r *http.Request, err error) {
	span := telemetry.CurrentSpan(r.Context())

	log.
		WithSpanAttrs(span,
			attribute.String("request.method", r.Method),
			attribute.String("request.url", r.URL.String()),
			attribute.String("trace", string(debug.Stack())),
		).
		Error(err.Error())
}

func (app *application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	span := telemetry.CurrentSpan(r.Context())
	span.SetAttributes(
		attribute.Int("response.status", status),
		attribute.String("response.message", message),
	)
	span.SetStatus(codes.Error, message)

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		span.RecordError(err)
		app.reportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	span := telemetry.CurrentSpan(r.Context())
	span.RecordError(err)
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

//lint:ignore U1000 useful templating code, remove lint:ignore when using for the first time
func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

//lint:ignore U1000 useful templating code, remove lint:ignore when using for the first time
func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
	}
}
