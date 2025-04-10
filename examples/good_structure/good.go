// Package goodstructure demonstrates proper Go package organization.
// Note that the package name matches the last component of the directory path.
package goodstructure

import (
	"fmt"
)

// Configuration represents application settings
type Configuration struct {
	DatabaseURL string
	APIKey      string
	Debug       bool
}

// defaultConfig is package-private (not exported)
var defaultConfig = Configuration{
	Debug: true,
}

// NewConfiguration creates a Configuration with defaults
func NewConfiguration() Configuration {
	return defaultConfig
}

// PrintInfo is a public function that prints configuration info
func PrintInfo(config Configuration) {
	fmt.Println("Configuration:")
	fmt.Printf("  Database: %s\n", maskSensitiveData(config.DatabaseURL))
	fmt.Printf("  Debug: %v\n", config.Debug)
}

// maskSensitiveData is package-private (not exported)
// This encapsulation keeps implementation details hidden
func maskSensitiveData(data string) string {
	if len(data) < 5 {
		return "***"
	}
	return data[:3] + "..." + data[len(data)-3:]
}

/*
ADVANTAGES OF THIS STRUCTURE:

1. Clear Package Boundaries:
   - Package name matches directory name
   - All related functionality is grouped together
   - Proper encapsulation with private functions (lowercase) and public APIs (uppercase)

2. Import Clarity:
   - Importing this package is clear: "go-examples/examples/good_structure"
   - Usage is intuitive: goodstructure.NewConfiguration()

3. Maintainability:
   - Adding new functionality is straightforward - just add more files to this directory
   - All implementation details stay encapsulated in the package
   - Clear separation between public API and private implementation

4. Testability:
   - Package boundaries create natural units for testing
   - Public APIs are clearly defined for test cases
*/

