package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes
	//a.use(a.Middleware.CheckRemember)

	// add routes here
	a.get("/", a.Handlers.Home)

	//svelte views routes
	//a.get("/svh/{module}", a.Handlers.SvelteViews)
	// pass all svelte subroutes back so that svelte can process them
	//a.get("/svh/{module}/*", a.Handlers.SvelteViews)

	//test cache api's with this page
	//a.get("/cache-test", a.Handlers.ShowCachePage)

	// Mount User login/out routes
	//a.App.Routes.Mount("/users", a.UsersRoutes())
	// Mount utility routes
	//a.App.Routes.Mount("/util", a.UtilityRoutes())
	// mount routes from celeritas
	//a.App.Routes.Mount("/celeritas", celeritas.Routes())
	// mount api routes - these are exempt from nosurf middleware
	// create api template with celeritas make cacheapi
	//a.App.Routes.Mount("/api", a.ApiRoutes())
	// Mount websocket/Pusher routes
	//a.App.Routes.Mount("/pusher", a.PusherRoutes())
	//a.get("/private-message", a.Handlers.WsSendPrivateMessage)
	//a.get("/broadcast-public-message", a.Handlers.WsSendPublicMessage)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}
