package logger_test

import (
	"bytes"
	"strings"
	"testing"

	"go-examples/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// TestLoggerInterface checks that our logger satisfies the Logger interface
func TestLoggerInterface(_ *testing.T) {
	// This is primarily a compile-time check
	var _ logger.Logger = &mockLogger{}
}

// TestNewDevelopment tests the development logger creation
func TestNewDevelopment(t *testing.T) {
	log, err := logger.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create development logger: %v", err)
	}
	if log == nil {
		t.Fatal("Expected non-nil logger")
	}
}

// TestNewProduction tests the production logger creation
func TestNewProduction(t *testing.T) {
	log, err := logger.NewProduction()
	if err != nil {
		t.Fatalf("Failed to create production logger: %v", err)
	}
	if log == nil {
		t.Fatal("Expected non-nil logger")
	}
}

// TestNewCustom tests custom logger creation
func TestNewCustom(t *testing.T) {
	// Test with development encoder
	devLog := logger.NewCustom(zapcore.DebugLevel, false)
	if devLog == nil {
		t.Fatal("Expected non-nil development logger")
	}
	
	// Test with production encoder
	prodLog := logger.NewCustom(zapcore.InfoLevel, true)
	if prodLog == nil {
		t.Fatal("Expected non-nil production logger")
	}
}

// TestLoggerImplementation tests the implementation using a custom test logger
func TestLoggerImplementation(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a test logger that writes to the buffer
	testLogger := createBufferedTestLogger(&buf, zapcore.DebugLevel)

	// Test normal logging methods
	tests := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name: "Debug method",
			logFunc: func() {
				testLogger.Debug("debug message")
			},
			contains: "debug message",
		},
		{
			name: "Info method",
			logFunc: func() {
				testLogger.Info("info message")
			},
			contains: "info message",
		},
		{
			name: "Warn method",
			logFunc: func() {
				testLogger.Warn("warn message")
			},
			contains: "warn message",
		},
		{
			name: "Error method",
			logFunc: func() {
				testLogger.Error("error message")
			},
			contains: "error message",
		},
		{
			name: "Debugf method",
			logFunc: func() {
				testLogger.Debugf("debug %s", "formatted")
			},
			contains: "debug formatted",
		},
		{
			name: "Infof method",
			logFunc: func() {
				testLogger.Infof("info %s", "formatted")
			},
			contains: "info formatted",
		},
		{
			name: "Warnf method",
			logFunc: func() {
				testLogger.Warnf("warn %s", "formatted")
			},
			contains: "warn formatted",
		},
		{
			name: "Errorf method",
			logFunc: func() {
				testLogger.Errorf("error %s", "formatted")
			},
			contains: "error formatted",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Clear the buffer before each test
			buf.Reset()
			
			// Execute the log function
			test.logFunc()
			
			// Check if the output contains the expected string
			if !strings.Contains(buf.String(), test.contains) {
				t.Errorf("Expected log output to contain %q, got: %s", test.contains, buf.String())
			}
		})
	}
}

// TestWithMethod tests the With method for context
func TestWithMethod(t *testing.T) {
	var buf bytes.Buffer
	testLogger := createBufferedTestLogger(&buf, zapcore.DebugLevel)
	
	// Create a derived logger with key-value context
	contextLogger := testLogger.With("key", "value")
	
	// Log a message with the context logger
	contextLogger.Info("message with context")
	
	// Check output contains both the message and context
	output := buf.String()
	if !strings.Contains(output, "message with context") {
		t.Errorf("Expected output to contain message, got: %s", output)
	}
	if !strings.Contains(output, "key") || !strings.Contains(output, "value") {
		t.Errorf("Expected output to contain context key/value, got: %s", output)
	}
}

// TestLogLevelFiltering tests that log level filtering works
func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	
	// Create a logger with Info level
	infoLogger := createBufferedTestLogger(&buf, zapcore.InfoLevel)
	
	// Debug should not appear
	buf.Reset()
	infoLogger.Debug("debug message")
	if buf.Len() > 0 {
		t.Errorf("Debug message should not appear with InfoLevel logger, got: %s", buf.String())
	}
	
	// Info should appear
	buf.Reset()
	infoLogger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Errorf("Info message should appear with InfoLevel logger, got: %s", buf.String())
	}
	
	// Create a logger with Error level
	buf.Reset()
	errorLogger := createBufferedTestLogger(&buf, zapcore.ErrorLevel)
	
	// Info should not appear
	errorLogger.Info("info should not appear")
	if buf.Len() > 0 {
		t.Errorf("Info message should not appear with ErrorLevel logger, got: %s", buf.String())
	}
	
	// Error should appear
	buf.Reset()
	errorLogger.Error("error should appear")
	if !strings.Contains(buf.String(), "error should appear") {
		t.Errorf("Error message should appear with ErrorLevel logger, got: %s", buf.String())
	}
}

// Helper function to create a logger that writes to a buffer
func createBufferedTestLogger(buf *bytes.Buffer, level zapcore.Level) logger.Logger {
	// Create a custom encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	
	// Create the test core with the buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(buf),
		level,
	)
	
	// Create the zap logger
	zapLogger := zap.New(core)
	
	// Need to convert to our logger interface
	// Easiest way is to use the public constructor with our test logger
	customLogger := &zapLoggerForTest{sugar: zapLogger.Sugar()}
	return customLogger
}

// Implement the logger interface for testing
type zapLoggerForTest struct {
	sugar *zap.SugaredLogger
}

func (l *zapLoggerForTest) Debug(args ...interface{})                   { l.sugar.Debug(args...) }
func (l *zapLoggerForTest) Info(args ...interface{})                    { l.sugar.Info(args...) }
func (l *zapLoggerForTest) Warn(args ...interface{})                    { l.sugar.Warn(args...) }
func (l *zapLoggerForTest) Error(args ...interface{})                   { l.sugar.Error(args...) }
func (l *zapLoggerForTest) Fatal(args ...interface{})                   { l.sugar.Fatal(args...) }
func (l *zapLoggerForTest) Debugf(template string, args ...interface{}) { l.sugar.Debugf(template, args...) }
func (l *zapLoggerForTest) Infof(template string, args ...interface{})  { l.sugar.Infof(template, args...) }
func (l *zapLoggerForTest) Warnf(template string, args ...interface{})  { l.sugar.Warnf(template, args...) }
func (l *zapLoggerForTest) Errorf(template string, args ...interface{}) { l.sugar.Errorf(template, args...) }
func (l *zapLoggerForTest) Fatalf(template string, args ...interface{}) { l.sugar.Fatalf(template, args...) }

func (l *zapLoggerForTest) With(args ...interface{}) logger.Logger {
	return &zapLoggerForTest{sugar: l.sugar.With(args...)}
}

// Example testing structured logging with zaptest
func TestStructuredLogging(t *testing.T) {
	// zaptest.NewLogger creates a logger that writes to the test's log output
	testLogger := zaptest.NewLogger(t)
	sugar := testLogger.Sugar()
	
	// Create a structured field
	sugar.Infow("structured log message",
		"string_key", "string value",
		"int_key", 123,
		"bool_key", true,
	)
	
	// This doesn't return anything we can assert on directly,
	// but it will be displayed in the test output
}

// mockLogger is a mock implementation of Logger for testing
type mockLogger struct{}

func (l *mockLogger) Debug(_ ...interface{})                   {}
func (l *mockLogger) Info(_ ...interface{})                    {}
func (l *mockLogger) Warn(_ ...interface{})                    {}
func (l *mockLogger) Error(_ ...interface{})                   {}
func (l *mockLogger) Fatal(_ ...interface{})                   {}
func (l *mockLogger) Debugf(_ string, _ ...interface{})        {}
func (l *mockLogger) Infof(_ string, _ ...interface{})         {}
func (l *mockLogger) Warnf(_ string, _ ...interface{})         {}
func (l *mockLogger) Errorf(_ string, _ ...interface{})        {}
func (l *mockLogger) Fatalf(_ string, _ ...interface{})        {}
func (l *mockLogger) With(_ ...interface{}) logger.Logger      { return l }