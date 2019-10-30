package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// APIRouter - The router instance serving the whole application for HTTP Requests
var APIRouter = mux.NewRouter()

func init() {
	APIRouter.HandleFunc("/auth", AuthHandler).Name("auth").Methods(http.MethodPost)
}
