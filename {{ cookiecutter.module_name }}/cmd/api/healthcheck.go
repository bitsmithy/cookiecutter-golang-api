package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"{{ cookiecutter.module_path }}/internal/build"
)

func (s *server) healthcheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := map[string]string{
		"environment": s.config.env,
	}

	for k, v := range build.Info() {
		info[k] = v
	}

	env := envelope{
		"status":      "available",
		"system_info": info,
	}

	err := s.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		s.serverErrorResponse(w, r, err)
	}
}
