package webutil

import (
	"net/http"
)

// Handler that calls the handler count times.
func MultiRequest(count int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < count; i++ {
			handler.ServeHTTP(w, r)
		}
	})
}
