package slogger_test

import (
	"bytes"
	"go-examples/pkg/slogger"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupTestHandler creates a test handler that writes logs to a buffer
func setupTestHandler(buf *bytes.Buffer) *slog.Logger {
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}

// TestFatalLogging tests the Fatal logging method (without actually exiting)
func TestFatalLogging(t *testing.T) {
	// Create a copy of os.Exit to restore it later
	origExit := slogger.OsExit
	defer func() { slogger.OsExit = origExit }()

	// Mock os.Exit to track if it was called
	var exitCalled bool
	slogger.OsExit = func(code int) {
		exitCalled = true
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
	}

	// Create a buffer to capture log output
	var buf bytes.Buffer
	origLogger := slog.Default()
	slog.SetDefault(setupTestHandler(&buf))
	defer slog.SetDefault(origLogger)

	// Call Fatal
	logger := slogger.Logger{}
	logger.Fatal("fatal error", "key", "value")

	// Check if os.Exit was called
	if !exitCalled {
		t.Error("os.Exit was not called")
	}

	// Verify log output
	output := buf.String()
	if !strings.Contains(output, "fatal error") {
		t.Errorf("expected log to contain 'fatal error', got: %s", output)
	}
	if !strings.Contains(output, "key") || !strings.Contains(output, "value") {
		t.Errorf("expected log to contain structured data, got: %s", output)
	}
}

// TestErrorLogging tests the Error logging method
func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	origLogger := slog.Default()
	slog.SetDefault(setupTestHandler(&buf))
	defer slog.SetDefault(origLogger)

	logger := slogger.Logger{}
	logger.Error("error message", "count", 42)

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("expected log to contain 'error message', got: %s", output)
	}
	if !strings.Contains(output, "count") || !strings.Contains(output, "42") {
		t.Errorf("expected log to contain structured data, got: %s", output)
	}
}

// TestInfoLogging tests the Info logging method
func TestInfoLogging(t *testing.T) {
	var buf bytes.Buffer
	origLogger := slog.Default()
	slog.SetDefault(setupTestHandler(&buf))
	defer slog.SetDefault(origLogger)

	logger := slogger.Logger{}
	logger.Info("info message", "flag", true)

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Errorf("expected log to contain 'info message', got: %s", output)
	}
	if !strings.Contains(output, "flag") || !strings.Contains(output, "true") {
		t.Errorf("expected log to contain structured data, got: %s", output)
	}
}

// TestInitLogging tests the initialization function
func TestInitLogging(_ *testing.T) {
	// Since InitLogging returns a zero value Logger struct (which is valid),
	// we just need to verify it doesn't panic
	_ = slogger.InitLogging()
	
	// If we get here, the test passes
	// Additional verification would be difficult without exposing implementation details
}

// TestResponseLogger tests the ResponseLogger struct and its methods
func TestResponseLogger(t *testing.T) {
	var buf bytes.Buffer
	origLogger := slog.Default()
	slog.SetDefault(setupTestHandler(&buf))
	defer slog.SetDefault(origLogger)
	
	// Create a logger and response logger
	logger := slogger.Logger{}
	respLogger := logger.NewResponseLogger("req-123")
	
	// Create a mock request
	req := httptest.NewRequest("GET", "/test", nil)
	
	// Test Response method
	buf.Reset()
	respLogger.Response(200, req, "action", "get_user")
	
	output := buf.String()
	if !strings.Contains(output, "req-123") {
		t.Errorf("expected log to contain request ID, got: %s", output)
	}
	if !strings.Contains(output, "code") || !strings.Contains(output, "200") {
		t.Errorf("expected log to contain status code, got: %s", output)
	}
	if !strings.Contains(output, "action") || !strings.Contains(output, "get_user") {
		t.Errorf("expected log to contain action, got: %s", output)
	}
}

// TestResponseErrorAndSend tests the ResponseErrorAndSend method
func TestResponseErrorAndSend(t *testing.T) {
	var buf bytes.Buffer
	origLogger := slog.Default()
	slog.SetDefault(setupTestHandler(&buf))
	defer slog.SetDefault(origLogger)
	
	// Create a logger and response logger
	logger := slogger.Logger{}
	respLogger := logger.NewResponseLogger("req-456")
	
	// Create a mock request and response writer
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	
	// Test ResponseErrorAndSend method
	respLogger.ResponseErrorAndSend(404, "Not Found", req, rec, "path", "/users/123")
	
	// Verify log output
	output := buf.String()
	if !strings.Contains(output, "req-456") {
		t.Errorf("expected log to contain request ID, got: %s", output)
	}
	if !strings.Contains(output, "message") || !strings.Contains(output, "Not Found") {
		t.Errorf("expected log to contain error message, got: %s", output)
	}
	
	// Verify HTTP response
	resp := rec.Result()
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}()
	
	if resp.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", resp.StatusCode)
	}
	
	// Read response body
	bodyBuf := new(bytes.Buffer)
	if _, err := bodyBuf.ReadFrom(resp.Body); err != nil {
		t.Errorf("error reading response body: %v", err)
	}
	body := bodyBuf.String()
	
	if !strings.Contains(body, "404 Not Found") {
		t.Errorf("expected response body to contain error message, got: %s", body)
	}
}