package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *server) healthcheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": s.config.env,
			"version":     version,
		},
	}

	err := s.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
