package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests and responses with default options
func LoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return LoggingMiddlewareWithConfig(logger, DefaultConfig())
}

// LoggingMiddlewareWithConfig logs HTTP requests and responses with custom configuration
func LoggingMiddlewareWithConfig(logger *Logger, config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract client IP
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			// Check if we should log this IP
			if !shouldLogIP(ip, config.FilterIPs) {
				next.ServeHTTP(w, r)
				return
			}

			requestLog := formatRequestLog(r)
			logger.Info("REQUEST: %s", requestLog)

			// Capture request body request logging is enabled
			var requestBody string
			if config.LogRequestBody && r.Body != nil && r.Method != http.MethodGet {
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Truncate if needed
				requestBody = string(bodyBytes)
				if len(requestBody) > config.MaxBodySize {
					requestBody = requestBody[:config.MaxBodySize] + "... [truncated]"
				}

				if requestBody != "" {
					logger.Info("REQUEST BODY: %s", requestBody)
				}
			}

			// Create custom response writer if response logging is enabled
			var rw *responseWriter
			if config.LogResponseBody {
				rw = &responseWriter{
					ResponseWriter: w,
					statusCode:     http.StatusOK,
					body:           &bytes.Buffer{},
				}
				next.ServeHTTP(rw, r)
			} else {
				next.ServeHTTP(w, r)
			}

			duration := time.Since(start)

			// Log request details
			logger.Info(
				"COMPLETED: [%s] (%s) URI: %s Address: %s",
				duration,
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
			)

			// Log response details if enabled
			if config.LogResponseBody && rw != nil {
				responseBody := rw.body.String()
				if len(responseBody) > config.MaxBodySize {
					responseBody = responseBody[:config.MaxBodySize] + "... [truncated]"
				}

				if responseBody != "" {
					logger.Info(
						"RESPONSE: Status: %d Body: %s",
						rw.statusCode,
						responseBody,
					)
				}
			}
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

// formatRequestLog formats the request details for logging
func formatRequestLog(r *http.Request) string {
	return fmt.Sprintf(
		"(%s) URI: %s Address: %s User-Agent: %s",
		r.Method,
		r.RequestURI,
		r.RemoteAddr,
		r.UserAgent(),
	)
}
