package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route - Struct used to define main aspects of an API route
type Route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Name    string
	Method  string
}

// Routes - Routes definition to be instantiated by the API
var Routes []Route = []Route{
	{"/auth", AuthHandler, "auth", http.MethodPost},
	{"/ws", WsHandler, "ws", http.MethodGet},
}

// Router - The router instance serving the whole application for HTTP Requests
var Router = mux.NewRouter()

func init() {
	for _, route := range Routes {
		Router.HandleFunc(route.Path, route.Handler).
			Name(route.Name).
			Methods(route.Method)
	}
	Router.Use(ValidateTokenMiddleware)
}
