package middleware

import (
	"os"
	"strconv"
)

// Config represents the logger middleware configuration
type Config struct {
	LogResponseBody bool
	LogRequestBody  bool
	FilterIPs       []string
	MaxBodySize     int
}

// DefaultConfig returns the default middleware configuration
func DefaultConfig() *Config {
	return &Config{
		LogResponseBody: true,
		LogRequestBody:  true,
		FilterIPs:       []string{}, // Empty means log every IP
		MaxBodySize:     1024,       // Truncate bodies larger than 1KB
	}
}

func CustomConfig() *Config {
	// Can't use my getter becuase of a cycle dependency
	logResponseBodyEnv := os.Getenv("LOG_RESPONSE_BODY")
	if logResponseBodyEnv == "" {
		logResponseBodyEnv = "true"
	}
	logResponseBody, err := strconv.ParseBool(logResponseBodyEnv)
	if err != nil {
		logResponseBody = true
	}

	logRequestBodyEnv := os.Getenv("LOG_REQUEST_BODY")
	if logRequestBodyEnv == "" {
		logRequestBodyEnv = "true"
	}
	logRequestBody, err := strconv.ParseBool(logResponseBodyEnv)
	if err != nil {
		logResponseBody = true
	}

	return &Config{
		LogResponseBody: logResponseBody,
		LogRequestBody:  logRequestBody,
		FilterIPs:       []string{}, // Empty means log every IP
		MaxBodySize:     1024,       // Truncate bodies larger than 1KB
	}
}
