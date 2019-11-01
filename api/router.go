package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// APIRouter - The router instance serving the whole application for HTTP Requests
var APIRouter = mux.NewRouter()

func init() {
	APIRouter.HandleFunc("/auth", AuthHandler).
		Name("auth").
		Methods(http.MethodPost)
	APIRouter.HandleFunc("/ws", WsHandler).
		Name("ws").
		Methods(http.MethodGet)
	APIRouter.HandleFunc("/nickname", NicknameHandler).
		Name("nickname").
		Methods(http.MethodPut)
	APIRouter.HandleFunc("/rooms", RoomsHandler).
		Name("rooms").
		Methods(http.MethodPost, http.MethodPut)
	APIRouter.HandleFunc("/users", UsersHandler).
		Name("users").
		Methods(http.MethodGet)
	APIRouter.HandleFunc("/exit", ExitHandler).
		Name("exit").
		Methods(http.MethodGet)
}
