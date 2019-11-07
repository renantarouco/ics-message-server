package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestValidateToken(t *testing.T) {
	t.Run("excluded-route", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("nickname", "test")
		req, err := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong code, wanted %v got %v", http.StatusCreated, status)
		}
	})
	t.Run("no-auth-header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("wrong code, wanted %v got %v", http.StatusUnauthorized, status)
		}
	})
	t.Run("header-no-token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("wrong code, wanted %v got %v", http.StatusUnauthorized, status)
		}
	})
	t.Run("valid-token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, validToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong code, wanted %v got %v", http.StatusOK, status)
		}
		if tokenStr := req.Context().Value("tokenStr"); tokenStr == nil {
			t.Errorf("context doesn't contain tokenStr")
		}
	})
	t.Run("invalid-token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Error(err)
		}
		addBearerAuthHeader(req, invalidToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("wrong code, wanted %v got %v", http.StatusUnauthorized, status)
		}
	})
}
