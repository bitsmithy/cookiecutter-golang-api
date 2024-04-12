package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/attribute"

	"{{ cookiecutter.module_path }}/internal/response"
	"{{ cookiecutter.module_path }}/internal/telemetry"
)

func (app *application) status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, span := telemetry.Trace(r.Context(), "handler.status")
	defer span.End()

	statusKey := "status"
	statusVal := "up"
	span.SetAttributes(attribute.String(statusKey, statusVal))

	err := response.JSON(w, http.StatusOK, map[string]string{
		statusKey: statusVal,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}
