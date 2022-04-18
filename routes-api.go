package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) ApiRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(mux chi.Router) {
		// add any api routes here
	})

	return r
}
