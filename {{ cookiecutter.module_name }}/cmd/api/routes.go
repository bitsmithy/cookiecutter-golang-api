package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) mux() *httprouter.Router {
	mux := httprouter.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.GET("/status", app.status)

	return mux
}

func (app *application) routes() http.Handler {
	return app.instrument(app.recoverPanic(app.logAccess(app.mux())))
}
