// Package main provides a simple interactive calculator application
// that uses the calculator package with logging capabilities.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-examples/pkg/calculator"
	"go-examples/pkg/logger"
)

func main() {
	// Initialize logger
	log, err := logger.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	log.Info("Starting calculator application")

	// Create calculator instance with logger
	calc := calculator.NewCalculator(log)
	fmt.Println("Simple Calculator")
	fmt.Println("=================")
	fmt.Println("Available operations: add, subtract, quit")
	fmt.Println("Example usage: add 5 3")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		log.Debugf("User input: %s", input)

		if input == "quit" || input == "exit" || input == "q" {
			log.Info("User requested to quit application")
			fmt.Println("Goodbye!")
			break
		}

		result, err := processCommand(input, calc, log)
		if err != nil {
			log.Warnf("Command processing error: %v", err)
			fmt.Printf("Error: %s\n", err)
			continue
		}

		log.Infof("Successful calculation, result: %d", result)
		fmt.Printf("Result: %d\n", result)
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("Scanner error: %v", err)
		fmt.Fprintf(os.Stderr, "Reading input: %s\n", err)
		os.Exit(1)
	}

	log.Info("Application shutting down")
}

func processCommand(input string, calc *calculator.Calculator, log logger.Logger) (int, error) {
	// Split the input into command and arguments
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid input, expected format: <operation> <number1> <number2>")
	}

	command := strings.ToLower(parts[0])

	// Parse the numbers
	a, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("first number is invalid: %v", err)
	}

	b, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("second number is invalid: %v", err)
	}

	// Perform the operation
	log.Debugf("Processing command: %s with arguments %d and %d", command, a, b)

	switch command {
	case "add":
		return calc.Add(a, b), nil
	case "subtract":
		return calc.Subtract(a, b), nil
	default:
		return 0, fmt.Errorf("unknown operation: %s, supported operations are 'add' and 'subtract'", command)
	}
}
