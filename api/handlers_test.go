package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func TestAuthHandler(t *testing.T) {
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/auth", strings.NewReader(""))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("wrong code, wanted %v got %v", http.StatusForbidden, status)
		}
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/auth", strings.NewReader("nickname=1test1"))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("wrong code, wanted %v got %v", http.StatusForbidden, status)
		}
	})
	t.Run("valid-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/auth", strings.NewReader("nickname=test1"))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong code, wanted %v got %v", http.StatusCreated, status)
		}
	})
	t.Run("duplicate-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/auth", strings.NewReader("nickname=test2"))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong code, wanted %v got %v", http.StatusCreated, status)
		}
		req, err = http.NewRequest("POST", "/auth", strings.NewReader("nickname=test2"))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr = httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("wrong code, wanted %v got %v", http.StatusForbidden, status)
		}
	})
	t.Run("valid-token", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/auth", strings.NewReader("nickname=test3"))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong code, wanted %v got %v", http.StatusCreated, status)
		}
		var bodyMap map[string]string
		err = json.Unmarshal(rr.Body.Bytes(), &bodyMap)
		if err != nil {
			t.Error(err)
		}
		tokenStr, ok := bodyMap["token"]
		if !ok {
			t.Error("body doesn't contain 'token' key")
		}
		token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			t.Error(err)
		}
		if !token.Valid {
			t.Error("invalid token received")
		}
	})
}
