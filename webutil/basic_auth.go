package webutil

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// Parses a request and returns the username and password provided. Returns two
// empty strings if the request does not have properly-formatted authorization.
func parseBasicAuthRequest(r *http.Request) (username, password string) {
	auth := r.Header.Get("Authorization")

	s := strings.SplitN(auth, " ", 2)
	if len(s) != 2 || s[0] != "Basic" {
		return "", ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", ""
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return "", ""
	}

	username = pair[0]
	password = pair[1]

	return
}

// Protects a handler with basic authentication.
func BasicAuthHandler(realm, username, password string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provided_username, provided_password := parseBasicAuthRequest(r)

		if provided_username != username || provided_password != password {
			log.Printf("Unauthenticated attempt at %s realm", realm)
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}
