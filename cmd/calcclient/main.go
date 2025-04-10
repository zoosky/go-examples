// Package main provides a CLI client for the calculator microservice
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Configuration holds client configuration
type Configuration struct {
	ServerURL string
	Timeout   time.Duration
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

	// Check if the service is available
	if !checkServiceHealth(config) {
		fmt.Println("Error: Calculator service is not available")
		os.Exit(1)
	}

	fmt.Println("Calculator Client")
	fmt.Println("================")
	fmt.Printf("Connected to: %s\n", config.ServerURL)
	fmt.Println("Available operations: add, subtract, multiply, divide, quit")
	fmt.Println("Example usage: add 5 3")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		fmt.Printf("Executing: %s\n", input)

		if input == "quit" || input == "exit" || input == "q" {
			fmt.Println("Goodbye!")
			break
		}

		result, err := processCommand(input, config)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		fmt.Printf("Result: %d\n", result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Reading input: %s\n", err)
		os.Exit(1)
	}
}

// parseFlags parses command line flags and returns configuration
func parseFlags() Configuration {
	serverURL := flag.String("server", "http://localhost:8080", "Calculator service URL")
	timeout := flag.Int("timeout", 5, "Request timeout in seconds")
	flag.Parse()

	return Configuration{
		ServerURL: *serverURL,
		Timeout:   time.Duration(*timeout) * time.Second,
	}
}

// checkServiceHealth verifies if the calculator service is available
func checkServiceHealth(config Configuration) bool {
	client := &http.Client{
		Timeout: config.Timeout,
	}

	resp, err := client.Get(fmt.Sprintf("%s/health", config.ServerURL))
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Health check failed with status: %s\n", resp.Status)
		return false
	}

	var healthResp map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		fmt.Printf("Failed to parse health response: %v\n", err)
		return false
	}

	return healthResp["status"]
}

// processCommand processes the user command and calls the API
func processCommand(input string, config Configuration) (int, error) {
	// Split the input into command and arguments
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid input, expected format: <operation> <number1> <number2>")
	}

	operation := strings.ToLower(parts[0])
	
	// Validate operation
	switch operation {
	case "add", "subtract", "multiply", "divide":
		// Valid operations
	default:
		return 0, fmt.Errorf("unknown operation: %s, supported operations are add, subtract, multiply, and divide", operation)
	}

	// Parse the numbers
	a, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("first number is invalid: %v", err)
	}

	b, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("second number is invalid: %v", err)
	}

	// Prepare the API request
	reqBody := CalculationRequest{
		Operation: operation,
		A:         a,
		B:         b,
	}

	return callCalculateAPI(reqBody, config)
}

// callCalculateAPI calls the calculate API endpoint
func callCalculateAPI(req CalculationRequest, config Configuration) (int, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/calculate", config.ServerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %v", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var calcResp CalculationResponse
	if err := json.Unmarshal(body, &calcResp); err != nil {
		return 0, fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for API errors
	if !calcResp.Success {
		return 0, fmt.Errorf("API error: %s", calcResp.Error)
	}

	return calcResp.Result, nil
}