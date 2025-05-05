package middleware

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Logger represents a custom logger with support for different log levels
type Logger struct {
	prefix    string
	debugMode bool
	output    io.Writer
}

// Speed up the traffic, unfortunately the logger comes at a cost of roughly x3 slowdown
func NewNoLogger() *Logger {
	return &Logger{
		prefix:    "",
		output:    io.Discard,
		debugMode: false,
	}
}

// New creates a new Logger instance
func New(out io.Writer, prefix string, debugMode bool) *Logger {
	return &Logger{
		prefix:    prefix,
		output:    out,
		debugMode: debugMode,
	}
}

// NewDefaultLogger creates a logger with default settings - writes to stdout with given prefix and act -
func NewDefaultLogger(prefix string) *Logger {
	debuggerOn, err := strconv.ParseBool(strings.ToLower(os.Getenv("LOGGER_DEBUG")))
	if err != nil {
		debuggerOn = false
	}
	debugMode := debuggerOn
	return New(os.Stdout, prefix, debugMode)
}

// FileDefaultLogger creates a new logger that writes to a file
func FileDefaultLogger(dir, filename, prefix string) (*Logger, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Join the directory and filename into a full path
	fullPath := filepath.Join(dir, filename)

	// Open or create the log file
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// If there's an error opening the file, we should handle file cleanup
	// to prevent file descriptor leaks

	debuggerOn := false
	debugStr := os.Getenv("LOGGER_DEBUG")
	if debugStr != "" {
		var err error
		debuggerOn, err = strconv.ParseBool(strings.ToLower(debugStr))
		if err != nil {
			// Keep default as false
			debuggerOn = false
		}
	}

	return New(file, prefix, debuggerOn), nil
}
