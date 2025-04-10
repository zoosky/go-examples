package calculator_test

import (
	"fmt"
	"testing"

	"go-examples/pkg/calculator"
	"go-examples/pkg/logger"
	"go.uber.org/zap/zapcore"
)

// setupTestLogger creates a logger suitable for tests
func setupTestLogger() logger.Logger {
	return logger.NewCustom(zapcore.DebugLevel, false)
}

func TestAdd(t *testing.T) {
	// Create test logger
	log := setupTestLogger()

	// Create calculator with test logger
	calc := calculator.NewCalculator(log)
	// Define test cases
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        5,
			b:        3,
			expected: 8,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: -5,
		},
		{
			name:     "mixed sign numbers",
			a:        5,
			b:        -3,
			expected: 2,
		},
		{
			name:     "zero values",
			a:        0,
			b:        0,
			expected: 0,
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := calc.Add(tc.a, tc.b)
			if got != tc.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	// Create test logger
	log := setupTestLogger()

	// Create calculator with test logger
	calc := calculator.NewCalculator(log)

	// Define test cases
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        5,
			b:        3,
			expected: 2,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: 1,
		},
		{
			name:     "mixed sign numbers",
			a:        5,
			b:        -3,
			expected: 8,
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := calc.Subtract(tc.a, tc.b)
			if got != tc.expected {
				t.Errorf("Subtract(%d, %d) = %d; want %d", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

// testOperation is a common testing function to reduce code duplication
func testOperation(t *testing.T, operation string, testCases []struct {
	name     string
	a, b     int
	expected int
}) {
	// Create test logger
	log := setupTestLogger()
	
	// Create calculator with test logger
	calc := calculator.NewCalculator(log)
	
	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got int
			
			// Call the appropriate operation based on the parameter
			switch operation {
			case "multiply":
				got = calc.Multiply(tc.a, tc.b)
			case "divide":
				got = calc.Divide(tc.a, tc.b)
			default:
				t.Fatalf("Unknown operation: %s", operation)
			}
			
			// Check the result
			if got != tc.expected {
				t.Errorf("%s(%d, %d) = %d; want %d", operation, tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        5,
			b:        3,
			expected: 15,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -3,
			expected: 6,
		},
		{
			name:     "mixed sign numbers",
			a:        5,
			b:        -3,
			expected: -15,
		},
		{
			name:     "multiply by zero",
			a:        5,
			b:        0,
			expected: 0,
		},
		{
			name:     "zero multiplied by number",
			a:        0,
			b:        5,
			expected: 0,
		},
		{
			name:     "large numbers",
			a:        1000,
			b:        1000,
			expected: 1000000,
		},
	}
	
	testOperation(t, "multiply", testCases)
}

func TestDivide(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        10,
			b:        2,
			expected: 5,
		},
		{
			name:     "negative numbers",
			a:        -10,
			b:        -2,
			expected: 5,
		},
		{
			name:     "mixed sign numbers",
			a:        10,
			b:        -2,
			expected: -5,
		},
		{
			name:     "integer division truncation",
			a:        11,
			b:        2,
			expected: 5,
		},
		{
			name:     "division by zero",
			a:        10,
			b:        0,
			expected: 0, // As per our implementation, division by zero returns 0
		},
		{
			name:     "zero divided by number",
			a:        0,
			b:        5,
			expected: 0,
		},
	}
	
	testOperation(t, "divide", testCases)
}

// Example functions are treated as documentation and also as tests.
// These examples appear in the generated documentation.
func ExampleAdd() {
	// Using the functional version for backward compatibility
	sum := calculator.Add(5, 3)
	fmt.Println(sum)
	// Output: 8
}

func ExampleSubtract() {
	// Using the functional version for backward compatibility
	difference := calculator.Subtract(5, 3)
	fmt.Println(difference)
	// Output: 2
}

// Examples using the object-oriented version with logger
func ExampleCalculator_Add() {
	// Create a development logger
	log, _ := logger.NewDevelopment()

	// Create a calculator with the logger
	calc := calculator.NewCalculator(log)

	// Perform calculation with logging
	sum := calc.Add(5, 3)
	fmt.Println(sum)
	// Output: 8
}

func ExampleCalculator_Subtract() {
	// Create a development logger
	log, _ := logger.NewDevelopment()

	// Create a calculator with the logger
	calc := calculator.NewCalculator(log)

	// Perform calculation with logging
	difference := calc.Subtract(5, 3)
	fmt.Println(difference)
	// Output: 2
}

func ExampleMultiply() {
	// Using the functional version for backward compatibility
	product := calculator.Multiply(5, 3)
	fmt.Println(product)
	// Output: 15
}

func ExampleCalculator_Multiply() {
	// Create a development logger
	log, _ := logger.NewDevelopment()
	
	// Create a calculator with the logger
	calc := calculator.NewCalculator(log)
	
	// Perform calculation with logging
	product := calc.Multiply(5, 3)
	fmt.Println(product)
	// Output: 15
}

func ExampleDivide() {
	// Using the functional version for backward compatibility
	quotient := calculator.Divide(10, 2)
	fmt.Println(quotient)
	// Output: 5
}

func ExampleCalculator_Divide() {
	// Create a development logger
	log, _ := logger.NewDevelopment()
	
	// Create a calculator with the logger
	calc := calculator.NewCalculator(log)
	
	// Perform calculation with logging
	quotient := calc.Divide(10, 2)
	fmt.Println(quotient)
	// Output: 5
}

// ----------------------
// Benchmark Tests
// ----------------------

// Basic operation benchmarks
func BenchmarkAdd(b *testing.B) {
	// Create a no-op logger to minimize logging overhead
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		calc.Add(5, 3)
	}
}

func BenchmarkSubtract(b *testing.B) {
	// Create a no-op logger to minimize logging overhead
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Subtract(5, 3)
	}
}

func BenchmarkMultiply(b *testing.B) {
	// Create a no-op logger to minimize logging overhead
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Multiply(5, 3)
	}
}

// Benchmarks with different input sizes
func BenchmarkMultiplySmall(b *testing.B) {
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Multiply(5, 3) // Small numbers
	}
}

func BenchmarkMultiplyMedium(b *testing.B) {
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Multiply(1000, 1000) // Medium numbers
	}
}

func BenchmarkMultiplyLarge(b *testing.B) {
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Multiply(1000000, 1000000) // Large numbers
	}
}

func BenchmarkDivide(b *testing.B) {
	// Create a no-op logger to minimize logging overhead
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Divide(10, 2)
	}
}

func BenchmarkDivideByZeroSafety(b *testing.B) {
	// Create a no-op logger to minimize logging overhead
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Divide(10, 0) // Tests the zero check
	}
}

// Benchmarks with different logger configurations
func BenchmarkAddWithRealLogger(b *testing.B) {
	// Use a development logger (with actual logging overhead)
	log, _ := logger.NewDevelopment()
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Add(5, 3)
	}
}

func BenchmarkAddWithNoLogger(b *testing.B) {
	// Use a no-op logger (minimal overhead)
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Add(5, 3)
	}
}

// Function-style vs method-style comparison
func BenchmarkAddFunction(b *testing.B) {
	// Using the package-level function
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculator.Add(5, 3)
	}
}

func BenchmarkAddMethod(b *testing.B) {
	// Using the method with a pre-initialized calculator
	log := noOpBenchLogger{}
	calc := calculator.NewCalculator(log)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Add(5, 3)
	}
}

// No-op logger implementation for benchmarks
type noOpBenchLogger struct{}

func (l noOpBenchLogger) Debug(_ ...interface{})              {}
func (l noOpBenchLogger) Info(_ ...interface{})               {}
func (l noOpBenchLogger) Warn(_ ...interface{})               {}
func (l noOpBenchLogger) Error(_ ...interface{})              {}
func (l noOpBenchLogger) Fatal(_ ...interface{})              {}
func (l noOpBenchLogger) Debugf(_ string, _ ...interface{})   {}
func (l noOpBenchLogger) Infof(_ string, _ ...interface{})    {}
func (l noOpBenchLogger) Warnf(_ string, _ ...interface{})    {}
func (l noOpBenchLogger) Errorf(_ string, _ ...interface{})   {}
func (l noOpBenchLogger) Fatalf(_ string, _ ...interface{})   {}
func (l noOpBenchLogger) With(_ ...interface{}) logger.Logger { return l }
