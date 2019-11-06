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

const (
	validToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzI3ODY0MDAsImlzcyI6Im5pY2szIiwic3ViIjoidW5hbWVkIn0.yJ-vyVz_5uJW4xyxNyA4TI00K98GfE4EkgYxBi8-w3c"
	invalidToken = "eyJhbGcikpXVCJ9.eyJpYXQiOjE1NzI3ODY0MDAsImlzcyI6Im5pY2szIiwic3ViIjoidW5hbWVkIn0.yJ-vyVz_5uJW4xyxNyA4TI00K98GfE4EkgYxBi8-w3c"
)

func todoTest(t *testing.T) {
	t.Error("this test is not yet implemented. DO IT!")
}

func testStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int) {
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("wrong code, wanted %v got %v", expectedStatusCode, status)
	}
}

func addBearerAuthHeader(r *http.Request, tokenStr string) {
	bearerToken := fmt.Sprintf("Bearer %s", "")
	r.Header.Add("Authorization", bearerToken)
}

func generateBasicNamingTest(fieldName, httpMethod, route, tokenStr string) func(t *testing.T, roomName string, expectedStatusCode int) *httptest.ResponseRecorder {
	return func(t *testing.T, name string, expectedStatusCode int) *httptest.ResponseRecorder {
		formData := url.Values{}
		formData.Set(fieldName, name)
		req, err := http.NewRequest(httpMethod, route, strings.NewReader(formData.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
		if tokenStr != "" {
			addBearerAuthHeader(req, tokenStr)
		}
		rr := httptest.NewRecorder()
		APIRouter.ServeHTTP(rr, req)
		testStatusCode(t, rr, expectedStatusCode)
		return rr
	}
}
