package api

import (
	"net/http"
	"strings"
)

// ValidateTokenMiddleware - Checks validity and extracts a token in the
// Authorization header
func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		excludeRoutes := []string{"/auth"}
		requestPath := r.URL.Path
		for _, route := range excludeRoutes {
			if requestPath == route {
				next.ServeHTTP(w, r)
				return
			}
		}
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		tokenStr = strings.TrimLeft(tokenStr, " ")
		err := IsTokenValid(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
