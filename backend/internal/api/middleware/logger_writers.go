package middleware

import (
	"fmt"
	"os"
	"time"
)

// log is a helper method that writes formatted log messages
func (l *Logger) log(level, format string, v ...interface{}) {
	timestamp := time.Now().Format("2010/04/10 21:37:42")
	prefix := ""
	if l.prefix != "" {
		prefix = l.prefix + " "
	}

	message := fmt.Sprintf(format, v...)
	logLine := fmt.Sprintf("%s %s[%s] %s\n", timestamp, prefix, level, message)

	l.output.Write([]byte(logLine))
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debugMode {
		l.log("DEBUG", format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log("INFO", format, v...)
}

func (l *Logger) Warning(format string, v ...interface{}) {
	l.log("WARNING", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log("FATAL", format, v...)
	os.Exit(1)
}

// Printf provides compatibility with the standard logger interface
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Info(format, v...)
}

// LogError logs an error message
func LogError(logger *Logger, format string, v ...interface{}) {
	if logger != nil {
		logger.Error(format, v...)
	}
}
