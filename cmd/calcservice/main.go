// Package main implements a calculator microservice
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-examples/pkg/calculator"
	"go-examples/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap/zapcore"
)

// Configuration holds all the server configuration
type Configuration struct {
	Port     int
	LogLevel string
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
	log, err := setupLogger(config.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log.Info("Starting calculator microservice")

	// Create calculator instance with logger
	calc := calculator.NewCalculator(log)

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
	flag.Parse()

	return Configuration{
		Port:     *port,
		LogLevel: *logLevel,
	}
}

// setupLogger creates and configures the logger
func setupLogger(level string) (logger.Logger, error) {
	var zapLevel zapcore.Level
	
	switch level {
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
}

// createCalculateHandler returns an HTTP handler for calculator operations
func createCalculateHandler(calc *calculator.Calculator, log logger.Logger) http.HandlerFunc {
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
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int, log logger.Logger) {
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