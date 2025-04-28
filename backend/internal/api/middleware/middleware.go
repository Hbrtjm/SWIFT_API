package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logs all HTTP requests, might replace it with external logger like prometheus or sentry + graphana 
func LogginMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			next.ServeHTTP(w, r)

			logger.Printf(
				"[%s] (%s) URI: %s Address: %s",
				time.Since(start),
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
			)
		})
	}
}

// ContentTypeMiddleware sets the Content-Type header for all responses
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}