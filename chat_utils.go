package main

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type chatToken struct {
	jwt.StandardClaims
	ClientIP string
	JoinTime int64
	RoomID   string
	UserID   string
}

var jwtKey = []byte("secret")

func enableCORS(router *mux.Router) http.Handler {
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

func validateNickname(nickname string) error {
	matched, err := regexp.MatchString("^[_a-zA-Z][_a-zA-Z-0-9]*", nickname)
	if !matched {
		return errors.New("invalid nickname")
	}
	if err != nil {
		return err
	}
	return nil
}
