package main

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type chatClaims struct {
	ClientIP string `json:"client_ip"`
	JoinTime int64  `json:"join_time"`
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
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

func validateToken(token string) (*chatClaims, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}
	claims := &chatClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
