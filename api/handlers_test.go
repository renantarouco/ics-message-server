package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJoinHandler(t *testing.T) {
	// subtests
	t.Run("no-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/join", strings.NewReader(""))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code, wanted %v got %v", http.StatusBadRequest, status)
		}
	})
	t.Run("valid-nickname", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/join", strings.NewReader("nickname=test"))
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
}
