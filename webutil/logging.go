package webutil

import (
	"log"
	"net/http"
	"time"
)

// A ResponseWriter proxy that keeps track of the status code.
type StatusResponseWriter struct {
	w http.ResponseWriter
	statusCode int
}

func (s *StatusResponseWriter) Header() http.Header {
	return s.w.Header()
}

func (s *StatusResponseWriter) Write(data []byte) (int, error) {
	return s.w.Write(data)
}

func (s *StatusResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.w.WriteHeader(statusCode)
}

func (s *StatusResponseWriter) StatusCode() int {
	if s.statusCode == 0 {
		// It's 0 if it's not been set, but then the http lib changes that to 200.
		return 200
	}
	return s.statusCode
}

func timeRequest(w *StatusResponseWriter, startTime time.Time, r *http.Request) {
	elapsed := time.Since(startTime)
	log.Printf("%s %s %d (%s)", r.Method, r.URL, w.StatusCode(), elapsed)
}

// Handler that logs request URLs and the time spent in them.
func LoggingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusW := StatusResponseWriter{w, 0}

		defer timeRequest(&statusW, time.Now(), r)

		handler.ServeHTTP(&statusW, r)
	})
}
