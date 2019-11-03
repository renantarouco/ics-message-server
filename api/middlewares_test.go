package api

import (
	"fmt"
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
		bearerToken := fmt.Sprintf("Bearer %s", "")
		req.Header.Add("Authorization", bearerToken)
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
		validToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzI3ODY0MDAsImlzcyI6Im5pY2szIiwic3ViIjoidW5hbWVkIn0.yJ-vyVz_5uJW4xyxNyA4TI00K98GfE4EkgYxBi8-w3c"
		bearerToken := fmt.Sprintf("Bearer %s", validToken)
		req.Header.Add("Authorization", bearerToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("wrong code, wanted %v got %v", http.StatusOK, status)
		}
	})
	t.Run("invalid-token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Error(err)
		}
		invalidToken := "eyJhbGcikpXVCJ9.eyJpYXQiOjE1NzI3ODY0MDAsImlzcyI6Im5pY2szIiwic3ViIjoidW5hbWVkIn0.yJ-vyVz_5uJW4xyxNyA4TI00K98GfE4EkgYxBi8-w3c"
		bearerToken := fmt.Sprintf("Bearer %s", invalidToken)
		req.Header.Add("Authorization", bearerToken)
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("wrong code, wanted %v got %v", http.StatusUnauthorized, status)
		}
	})
}
