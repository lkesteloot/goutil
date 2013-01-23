package webutil

import (
	"log"
	"net/http"
)

// Handler that rejects (with 405) any request with a method other than GET.
func GetHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			log.Printf("Method %s not allowed", r.Method)
			http.Error(w,
				http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

// Handler that accepts GET and POST requests, sending them to two different handlers,
// and rejects (with 405) any other method.
func GetPostHandler(getHandler, postHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getHandler.ServeHTTP(w, r)
		} else if r.Method == "POST" {
			postHandler.ServeHTTP(w, r)
		} else {
			log.Printf("Method %s not allowed", r.Method)
			http.Error(w,
				http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
		}
	})
}
