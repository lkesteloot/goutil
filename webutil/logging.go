package webutil

import (
	"log"
	"net/http"
	"time"
)

func timeRequest(startTime time.Time, r *http.Request) {
	elapsed := time.Since(startTime)
	log.Printf("%s %s (%s)", r.Method, r.URL, elapsed)
	// Unfortunately I don't see a way to get the status code from the response here.
}

// Handler that logs request URLs and the time spent in them.
func LoggingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer timeRequest(time.Now(), r)

		handler.ServeHTTP(w, r)
	})
}
