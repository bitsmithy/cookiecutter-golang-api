package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(s.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(s.methodNotAllowedResponse)

	router.GET("/healthcheck", s.healthcheckHandler)

	return s.recoverPanic(s.logRequest(router))
}
