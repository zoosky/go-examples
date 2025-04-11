// Package main implements a calculator microservice
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-examples/pkg/calculator"
	"go-examples/pkg/logger"
	"go-examples/pkg/slogger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap/zapcore"
)

// LoggerInterface defines a common interface for both logging systems
type LoggerInterface interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
	
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// SlogAdapter adapts the slogger to our common interface
type SlogAdapter struct {
	logger slogger.Logger
}

// Info logs an informational message
func (s *SlogAdapter) Info(args ...interface{}) {
	// Convert args to a message string and any key-value pairs
	if len(args) > 0 {
		msg, ok := args[0].(string)
		if ok {
			s.logger.Info(msg, args[1:]...)
		} else {
			s.logger.Info("info", args...)
		}
	}
}

// Error logs an error message
func (s *SlogAdapter) Error(args ...interface{}) {
	if len(args) > 0 {
		msg, ok := args[0].(string)
		if ok {
			s.logger.Error(msg, args[1:]...)
		} else {
			s.logger.Error("error", args...)
		}
	}
}

// Debug logs a debug message (maps to Info in slogger)
func (s *SlogAdapter) Debug(args ...interface{}) {
	// slogger doesn't have a Debug method, so use Info
	s.Info(args...)
}

// Warn logs a warning message (maps to Info in slogger)
func (s *SlogAdapter) Warn(args ...interface{}) {
	// slogger doesn't have a Warn method, so use Info
	s.Info(args...)
}

// Fatal logs a fatal error message and exits the program
func (s *SlogAdapter) Fatal(args ...interface{}) {
	if len(args) > 0 {
		msg, ok := args[0].(string)
		if ok {
			s.logger.Fatal(msg, args[1:]...)
		} else {
			s.logger.Fatal("fatal", args...)
		}
	}
}

// Infof logs an informational message with formatting
func (s *SlogAdapter) Infof(template string, args ...interface{}) {
	// slogger doesn't have formatted methods, so we'll format it ourselves
	s.logger.Info(fmt.Sprintf(template, args...))
}

// Errorf logs an error message with formatting
func (s *SlogAdapter) Errorf(template string, args ...interface{}) {
	s.logger.Error(fmt.Sprintf(template, args...))
}

// Warnf logs a warning message with formatting
func (s *SlogAdapter) Warnf(template string, args ...interface{}) {
	// slogger doesn't have Warn, so use Info
	s.logger.Info("WARN: " + fmt.Sprintf(template, args...))
}

// Fatalf logs a fatal error message with formatting and exits the program
func (s *SlogAdapter) Fatalf(template string, args ...interface{}) {
	s.logger.Fatal(fmt.Sprintf(template, args...))
}

// calculatorLoggerAdapter adapts our common interface to the calculator's logger interface
type calculatorLoggerAdapter struct {
	log LoggerInterface
}

func (a *calculatorLoggerAdapter) Debug(args ...interface{})              { a.log.Debug(args...) }
func (a *calculatorLoggerAdapter) Info(args ...interface{})               { a.log.Info(args...) }
func (a *calculatorLoggerAdapter) Warn(args ...interface{})               { a.log.Warn(args...) }
func (a *calculatorLoggerAdapter) Error(args ...interface{})              { a.log.Error(args...) }
func (a *calculatorLoggerAdapter) Fatal(args ...interface{})              { a.log.Fatal(args...) }
func (a *calculatorLoggerAdapter) Debugf(template string, args ...interface{})   { a.log.Infof(template, args...) }
func (a *calculatorLoggerAdapter) Infof(template string, args ...interface{})    { a.log.Infof(template, args...) }
func (a *calculatorLoggerAdapter) Warnf(template string, args ...interface{})    { a.log.Infof(template, args...) }
func (a *calculatorLoggerAdapter) Errorf(template string, args ...interface{})   { a.log.Errorf(template, args...) }
func (a *calculatorLoggerAdapter) Fatalf(template string, args ...interface{})   { a.log.Fatal(fmt.Sprintf(template, args...)) }
func (a *calculatorLoggerAdapter) With(_ ...interface{}) logger.Logger { return a }

// Configuration holds all the server configuration
type Configuration struct {
	Port      int
	LogLevel  string
	LogSystem string // "zap" or "slog"
}

// CalculationRequest represents a calculation API request
type CalculationRequest struct {
	Operation string `json:"operation"`
	A         int    `json:"a"`
	B         int    `json:"b"`
}

// CalculationResponse represents a calculation API response
type CalculationResponse struct {
	Result  int    `json:"result"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	// Parse configuration from command line flags
	config := parseFlags()

	// Initialize logger
	log, err := setupLogger(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log.Info("Starting calculator microservice")
	log.Infof("Using %s logging system", config.LogSystem)

	// Create calculator instance with logger
	var calcLogger logger.Logger
	if zapLogger, ok := log.(logger.Logger); ok {
		// If it's the original logger, use it directly
		calcLogger = zapLogger
	} else {
		// If it's the slog adapter, create a simple adapter for the calculator
		// The calculator expects the original logger interface
		calcLogger = &calculatorLoggerAdapter{log: log}
	}
	calc := calculator.NewCalculator(calcLogger)

	// Set up API routes
	router := mux.NewRouter()
	router.HandleFunc("/calculate", createCalculateHandler(calc, log)).Methods("POST")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Start server
	serverAddr := fmt.Sprintf(":%d", config.Port)
	log.Infof("Server starting on %s", serverAddr)
	
	// Create a server with graceful shutdown and security settings
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second, // Prevent Slowloris attacks
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-stop
	log.Info("Shutting down server...")
	log.Info("Server stopped")
}

// parseFlags parses command line flags and returns configuration
func parseFlags() Configuration {
	port := flag.Int("port", 8080, "Server port")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	logSystem := flag.String("log-system", "zap", "Logging system to use (zap or slog)")
	flag.Parse()

	return Configuration{
		Port:      *port,
		LogLevel:  *logLevel,
		LogSystem: strings.ToLower(*logSystem),
	}
}

// setupLogger creates and configures the logger based on the configuration
func setupLogger(config Configuration) (LoggerInterface, error) {
	switch config.LogSystem {
	case "slog":
		// Initialize structured logger (slogger)
		slog := slogger.InitLogging()
		return &SlogAdapter{logger: slog}, nil
		
	case "zap", "":
		// Initialize zap logger (original logger)
		var zapLevel zapcore.Level
		
		switch config.LogLevel {
		case "debug":
			zapLevel = zapcore.DebugLevel
		case "info":
			zapLevel = zapcore.InfoLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		default:
			zapLevel = zapcore.InfoLevel
		}
		
		// Using NewCustom which doesn't return error
		return logger.NewCustom(zapLevel, true), nil
		
	default:
		return nil, fmt.Errorf("unknown log system: %s, supported systems are 'zap' and 'slog'", config.LogSystem)
	}
}

// createCalculateHandler returns an HTTP handler for calculator operations
func createCalculateHandler(calc *calculator.Calculator, log LoggerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var req CalculationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendErrorResponse(w, "Invalid request format", http.StatusBadRequest, log)
			return
		}

		log.Infof("Calculation request: %+v", req)

		// Process calculation
		var result int

		switch req.Operation {
		case "add":
			result = calc.Add(req.A, req.B)
		case "subtract":
			result = calc.Subtract(req.A, req.B)
		case "multiply":
			result = calc.Multiply(req.A, req.B)
		case "divide":
			if req.B == 0 {
				sendErrorResponse(w, "Division by zero", http.StatusBadRequest, log)
				return
			}
			result = calc.Divide(req.A, req.B)
		default:
			sendErrorResponse(w, "Unknown operation: "+req.Operation, http.StatusBadRequest, log)
			return
		}

		// Send successful response
		resp := CalculationResponse{
			Result:  result,
			Success: true,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Errorf("Failed to encode response: %v", err)
		}
	}
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"status": true}); err != nil {
		// This would rarely happen, but we should handle it
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// sendErrorResponse sends an error response with the given message and status code
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int, log LoggerInterface) {
	log.Warnf("Error response: %s (code: %d)", message, statusCode)
	resp := CalculationResponse{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf("Failed to encode error response: %v", err)
		// In case we can't encode the JSON response, send a plain text error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}