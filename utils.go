package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// EnableCORS - Necessary headers to allow CORS
func EnableCORS(router *mux.Router) http.Handler {
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Content-Type",
		"Accept",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	})
	return handlers.CORS(headersOk, originsOk, methodsOk)(router)
}
