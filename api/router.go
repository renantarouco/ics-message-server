package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// APIRouter - The router instance serving the whole application for HTTP Requests
var APIRouter = mux.NewRouter()

func init() {
	APIRouter.HandleFunc("/join", JoinHandler).Name("join").Methods(http.MethodPost)
}

func init() {}
