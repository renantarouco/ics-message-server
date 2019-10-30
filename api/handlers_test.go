package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJoinHandler(t *testing.T) {
	// setup
	testNicknameFunc := func(t *testing.T, formBody string, expectedStatusCode int) {
		req, err := http.NewRequest(http.MethodPost, "/join", strings.NewReader(formBody))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		if status := rr.Code; status != expectedStatusCode {
			t.Errorf("wrong code, wanted %v got %v", expectedStatusCode, status)
		}
	}
	// subtests
	t.Run("no-nickname", func(t *testing.T) {
		testNicknameFunc(t, "", http.StatusBadRequest)
	})
	t.Run("valid-nickname", func(t *testing.T) {
		testNicknameFunc(t, "nickname=valid-nick", http.StatusCreated)
	})
	t.Run("invalid-nickname", func(t *testing.T) {
		testNicknameFunc(t, "nickname=(inváli$d", http.StatusBadRequest)
	})
}
