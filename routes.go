package main

import (
	"net/http"

	"github.com/brucebotes/celeritas"
	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes

	// add routes here
	a.get("/", a.Handlers.Home)
	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	// mount routes from celeritas
	a.App.Routes.Mount("/celeritas", celeritas.Routes())
	// mount api routes - these are exempt from nosurf middleware
	a.App.Routes.Mount("/api", a.ApiRoutes())

	return a.App.Routes
}
