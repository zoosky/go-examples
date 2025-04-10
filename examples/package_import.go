// Package examples demonstrates how Go's package structure works.
// This example shows the relationship between directory structure and package imports.
package examples

import (
	"fmt"
	
	// Import packages by their full module path + directory path
	"go-examples/pkg/calculator" // Imports the calculator package from pkg/calculator directory
	"go-examples/pkg/logger"     // Imports the logger package from pkg/logger directory
)

// PackageExample demonstrates Go's package structure 
func PackageExample() {
	// Initialize a logger
	log, err := logger.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	
	// Create a calculator with the logger
	calc := calculator.NewCalculator(log)
	
	// Use the calculator
	result := calc.Add(10, 20)
	
	// When we import "go-examples/pkg/calculator", we're referring to the package
	// defined in the directory "pkg/calculator". All files in that directory 
	// must belong to the same package named "calculator".
	fmt.Printf("Result: %d\n", result)
	
	// Why this structure is maintainable:
	fmt.Println("\nWhy Go's package structure is maintainable:")
	fmt.Println("1. Clear 1:1 mapping between directory names and package names")
	fmt.Println("2. Easy to locate code - if you want calculator functionality, look in pkg/calculator")
	fmt.Println("3. Packages are self-contained units - all calculator code is in one place")
	fmt.Println("4. Prevents circular dependencies due to package boundaries")
	fmt.Println("5. Allows for better organization of larger codebases")
	fmt.Println("6. Import paths clearly indicate where code comes from")
}

/*
Go Package Structure Best Practices:

1. Package == Directory
   - Each directory corresponds to exactly one package
   - All files in a directory must belong to the same package
   - The package name is typically the same as the directory name

2. Import Path vs Package Name
   - Import path: full path from module root (e.g., "go-examples/pkg/calculator")
   - Package name: just the name used in code (e.g., "calculator")
   
3. Standard Project Layout
   /cmd           - Main applications
   /pkg           - Library code that can be used by external applications
   /internal      - Private library code
   /api           - API definitions (OpenAPI/Swagger, protocol buffers, etc)
   /web           - Web assets
   /configs       - Configuration files
   /test          - Additional test applications and test data
   
4. Advantages
   - Easier navigation - directory structure maps to logical components
   - Better encapsulation - package boundaries enforce clean interfaces
   - Improved maintainability - related code stays together
   - Clearer dependencies - explicit import paths show relationships
*/

