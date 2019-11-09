package http

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// IsTokenValid - Checks if a token string is valid
func IsTokenValid(tokenStr string) error {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

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
