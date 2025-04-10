// Package badstructure demonstrates problematic Go package organization.
// This example shows what NOT to do when organizing Go code.
package badstructure

import (
	"fmt"
)

// This would typically be in a separate file but the same package
// We're putting everything in one file to demonstrate problems

// Problem 1: Mixed Concerns
// The same package contains unrelated functionalities (DB and HTTP utils)

// Database represents a database connection
type Database struct {
	URL      string
	Username string
	Password string
}

// HTTPClient represents an HTTP client
type HTTPClient struct {
	BaseURL    string
	UserAgent  string
	MaxRetries int
}

// Problem 2: Name conflicts with methods
// Both types have "Connect" methods, forcing awkward usage patterns

// Connect establishes a database connection
func (db Database) Connect() error {
	fmt.Printf("Connecting to database at %s\n", db.URL)
	return nil
}

// Connect establishes an HTTP connection
// This compiles but leads to confusing usage:
// - db.Connect() and client.Connect() do very different things
// - No namespace separation despite different purposes
func (client HTTPClient) Connect() error {
	fmt.Printf("Connecting to HTTP service at %s\n", client.BaseURL)
	return nil
}

// Problem 3: No encapsulation between components
// The validate function is reused inappropriately across components

func validate(s string) bool {
	return len(s) > 0
}

// ValidateDatabase checks database settings
func ValidateDatabase(db Database) bool {
	// Reuses the same validation logic for completely different types
	return validate(db.URL) && validate(db.Username)
}

// ValidateHTTPClient checks HTTP client settings
func ValidateHTTPClient(client HTTPClient) bool {
	// Reuses the same validation logic intended for DB validation
	return validate(client.BaseURL) && validate(client.UserAgent)
}

/*
PROBLEMS WITH THIS STRUCTURE:

1. Poor Separation of Concerns:
   - Database and HTTP functionality are mixed in one package
   - No logical boundaries between different components
   - Difficult to understand which parts relate to each other

2. Name Conflicts:
   - Methods with the same name on different types cause confusion
   - No namespacing to differentiate between components
   - Users must understand implementation details to know what each method does

3. Inappropriate Code Reuse:
   - Shared functions used across unrelated components
   - Changes to shared functions affect multiple systems unintentionally

4. Import Problems:
   - Importing this package pulls in ALL functionality, even if just one part is needed
   - No way to import just the database or just the HTTP components

5. Maintenance Nightmare:
   - Adding new functionality may require checking all existing code for conflicts
   - Changing "validate" affects both database and HTTP validation
   - No clear ownership of specific pieces of code

Better Approach:
   - Split into separate packages: "database" and "http"
   - Each package would be in its own directory with clear responsibility
   - Methods can have the same names without conflicts: database.Connect() vs http.Connect()
*/

