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
	t.Run("valid-token", func(t *testing.T) {

	})
}
