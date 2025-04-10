// Package slogger provides a wrapper around Go's structured logging (slog) package
// with additional features for HTTP response logging.
package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// Logger is a wrapper around slog that provides simpler methods for common logging levels.
type Logger struct{}

// OsExit is a variable that points to os.Exit to allow for testing
// without actually exiting the program.
var OsExit = os.Exit

// Fatal logs a message at fatal level and then exits the program with status code 1.
func (l *Logger) Fatal(msg string, args ...any) {
	slog.Log(context.Background(), slog.LevelError, msg, args...)
	OsExit(1)
}

// Error logs a message at error level.
func (l *Logger) Error(msg string, args ...any) {
	slog.Log(context.Background(), slog.LevelError, msg, args...)
}

// Info logs a message at info level.
func (l *Logger) Info(msg string, args ...any) {
	slog.Log(context.Background(), slog.LevelInfo, msg, args...)
}

// InitLogging initializes the structured logger with DEBUG level
// and returns a new Logger instance.
func InitLogging() Logger {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return Logger{}
}

// ResponseLogger provides logging utilities specifically for HTTP responses
// with request context information included.
type ResponseLogger struct {
	requestID string // Unique ID for the request
	logger    *Logger
}

// Response logs information about an HTTP response including status code and URI.
func (l *ResponseLogger) Response(code int, r *http.Request, args ...any) {
	params := append([]any{"code", code, "uri", r.RequestURI}, args...)
	l.logger.Info(l.requestID, params...)
}

// ResponseErrorAndSend logs an error response and sends it to the client.
func (l *ResponseLogger) ResponseErrorAndSend(code int, msg string, r *http.Request, w http.ResponseWriter, args ...any) {
	l.Response(code, r, append([]any{"message", msg}, args...)...)
	http.Error(w, fmt.Sprintf("%d %s", code, msg), code)
}

// NewResponseLogger creates a new ResponseLogger with the specified request ID.
func (l *Logger) NewResponseLogger(requestID string) *ResponseLogger {
	return &ResponseLogger{
		requestID: requestID,
		logger:    l,
	}
}
