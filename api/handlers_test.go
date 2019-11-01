package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func TestAuthHandler(t *testing.T) {
	nicknameTest := func(t *testing.T, httpMethod, route, nickname string, expectedStatusCode int) *httptest.ResponseRecorder {
		formData := url.Values{}
		formData.Set("nickname", nickname)
		req, err := http.NewRequest(httpMethod, route, strings.NewReader(formData.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != expectedStatusCode {
			t.Errorf("wrong code, wanted %v got %v", expectedStatusCode, status)
		}
		return rr
	}
	t.Run("no-nickname", func(t *testing.T) {
		nicknameTest(t, http.MethodPost, "/auth", "", http.StatusForbidden)
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		nicknameTest(t, http.MethodPost, "/auth", "1nick1", http.StatusForbidden)
	})
	t.Run("valid-nickname", func(t *testing.T) {
		nicknameTest(t, http.MethodPost, "/auth", "nick1", http.StatusCreated)
	})
	t.Run("duplicate-nickname", func(t *testing.T) {
		nicknameTest(t, http.MethodPost, "/auth", "nick2", http.StatusCreated)
		nicknameTest(t, http.MethodPost, "/auth", "nick2", http.StatusForbidden)
	})
	t.Run("valid-token", func(t *testing.T) {
		rr := nicknameTest(t, http.MethodPost, "/auth", "nick3", http.StatusCreated)
		var bodyMap map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &bodyMap)
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

func TestWsHandler(t *testing.T) {}

func TestNicknameHandler(t *testing.T) {}

func TestRoomsHandler(t *testing.T) {}

func TestUsersHandler(t *testing.T) {}

func TestExitHandler(t *testing.T) {}
