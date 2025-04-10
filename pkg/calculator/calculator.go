// Package calculator provides arithmetic operations on numbers.
package calculator

import (
	"go-examples/pkg/logger"
)

// Calculator provides arithmetic operations with logging capabilities
type Calculator struct {
	log logger.Logger
}

// NewCalculator creates a new Calculator instance with the provided logger
func NewCalculator(log logger.Logger) *Calculator {
	return &Calculator{
		log: log,
	}
}

// Add returns the sum of two integers.
// It's a simple function to demonstrate Go package functionality.
func (c *Calculator) Add(a, b int) int {
	c.log.Infof("Calculating addition: %d + %d", a, b)
	result := a + b
	c.log.Debugf("Addition result: %d", result)
	return result
}

// Subtract returns the difference between two integers.
// It subtracts the second argument from the first.
func (c *Calculator) Subtract(a, b int) int {
	c.log.Infof("Calculating subtraction: %d - %d", a, b)
	result := a - b
	c.log.Debugf("Subtraction result: %d", result)
	return result
}

// Multiply returns the product of two integers.
// It multiplies the first argument by the second.
func (c *Calculator) Multiply(a, b int) int {
	c.log.Infof("Calculating multiplication: %d * %d", a, b)
	result := a * b
	c.log.Debugf("Multiplication result: %d", result)
	return result
}

// Divide returns the quotient of two integers.
// It divides the first argument by the second.
func (c *Calculator) Divide(a, b int) int {
	c.log.Infof("Calculating division: %d / %d", a, b)
	if b == 0 {
		c.log.Error("Division by zero")
		return 0
	}
	result := a / b
	c.log.Debugf("Division result: %d", result)
	return result
}

// For backward compatibility with existing code, keep the original functions
// but they now use a default no-op logger

// Add returns the sum of two integers.
func Add(a, b int) int {
	// Create a calculator with a no-op logger for backward compatibility
	calc := NewCalculator(noOpLogger{})
	return calc.Add(a, b)
}

// Subtract returns the difference between two integers.
func Subtract(a, b int) int {
	// Create a calculator with a no-op logger for backward compatibility
	calc := NewCalculator(noOpLogger{})
	return calc.Subtract(a, b)
}

// Multiply returns the product of two integers.
func Multiply(a, b int) int {
	// Create a calculator with a no-op logger for backward compatibility
	calc := NewCalculator(noOpLogger{})
	return calc.Multiply(a, b)
}

// Divide returns the quotient of two integers.
func Divide(a, b int) int {
	// Create a calculator with a no-op logger for backward compatibility
	calc := NewCalculator(noOpLogger{})
	return calc.Divide(a, b)
}

// noOpLogger is a no-operation logger for backward compatibility
type noOpLogger struct{}

func (l noOpLogger) Debug(_ ...interface{})              {}
func (l noOpLogger) Info(_ ...interface{})               {}
func (l noOpLogger) Warn(_ ...interface{})               {}
func (l noOpLogger) Error(_ ...interface{})              {}
func (l noOpLogger) Fatal(_ ...interface{})              {}
func (l noOpLogger) Debugf(_ string, _ ...interface{})   {}
func (l noOpLogger) Infof(_ string, _ ...interface{})    {}
func (l noOpLogger) Warnf(_ string, _ ...interface{})    {}
func (l noOpLogger) Errorf(_ string, _ ...interface{})   {}
func (l noOpLogger) Fatalf(_ string, _ ...interface{})   {}
func (l noOpLogger) With(_ ...interface{}) logger.Logger { return l }
