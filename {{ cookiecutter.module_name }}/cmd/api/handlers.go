package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"{{ cookiecutter.module_path }}/internal/response"
)

func (app *application) status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
