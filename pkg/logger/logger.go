// Package logger provides a centralized logging configuration
// based on uber-go/zap.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger interface
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	With(args ...interface{}) Logger
}

// zapLogger wraps zap.SugaredLogger to implement our Logger interface
type zapLogger struct {
	sugar *zap.SugaredLogger
}

// NewDevelopment creates a logger with development-friendly defaults
func NewDevelopment() (Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return &zapLogger{sugar: sugar}, nil
}

// NewProduction creates a logger with production-friendly defaults
func NewProduction() (Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return &zapLogger{sugar: sugar}, nil
}

// NewCustom creates a logger with custom configuration
func NewCustom(level zapcore.Level, isProduction bool) Logger {
	// Create encoder config based on environment
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use JSON encoder for production, console encoder for development
	var encoder zapcore.Encoder
	if isProduction {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &zapLogger{sugar: logger.Sugar()}
}

// Implementation of Logger interface methods
func (l *zapLogger) Debug(args ...interface{})                   { l.sugar.Debug(args...) }
func (l *zapLogger) Info(args ...interface{})                    { l.sugar.Info(args...) }
func (l *zapLogger) Warn(args ...interface{})                    { l.sugar.Warn(args...) }
func (l *zapLogger) Error(args ...interface{})                   { l.sugar.Error(args...) }
func (l *zapLogger) Fatal(args ...interface{})                   { l.sugar.Fatal(args...) }
func (l *zapLogger) Debugf(template string, args ...interface{}) { l.sugar.Debugf(template, args...) }
func (l *zapLogger) Infof(template string, args ...interface{})  { l.sugar.Infof(template, args...) }
func (l *zapLogger) Warnf(template string, args ...interface{})  { l.sugar.Warnf(template, args...) }
func (l *zapLogger) Errorf(template string, args ...interface{}) { l.sugar.Errorf(template, args...) }
func (l *zapLogger) Fatalf(template string, args ...interface{}) { l.sugar.Fatalf(template, args...) }

func (l *zapLogger) With(args ...interface{}) Logger {
	return &zapLogger{sugar: l.sugar.With(args...)}
}
