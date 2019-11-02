package api

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
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
