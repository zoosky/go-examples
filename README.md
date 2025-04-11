# Go Examples

A collection of Go code examples demonstrating various aspects of Go development, including:

- Package structure
- CLI application development
- Microservice development
- RESTful API client
- Testing
- Logging
- Error handling

## Project Components

### 1. Calculator Package

- Located in: `pkg/calculator`
- Provides basic arithmetic operations: add, subtract, multiply, divide
- Includes testing and benchmarking examples
- Uses structured logging

### 2. Logger Package

- Located in: `pkg/logger`
- Wrapper around zap logging library
- Provides consistent logging interface across applications

### 3. SLogger Package

- Located in: `pkg/slogger`
- Wrapper around Go's standard library slog package
- Provides simplified structured logging interface

### 4. CLI Calculator App

- Located in: `cmd/app`
- Command-line application for basic arithmetic
- Uses the calculator package directly
- Interactive interface

### 5. Calculator Microservice

- Located in: `cmd/calcservice`
- RESTful API for calculator operations
- JSON request/response format
- Configurable port and log level
- Configurable logging system (ZAP or SLOG)
- Health check endpoint
- Graceful shutdown

### 6. Calculator API Client

- Located in: `cmd/calcclient`
- Command-line client that uses the calculator microservice
- Makes HTTP requests to the service
- Interactive interface
- Configurable server URL and timeout

## Getting Started

### Prerequisites

- Go 1.18 or later
- Make (for build automation)

### Building and Running

```bash
# Build all components
make build

# Run the direct CLI calculator
make run

# Run the calculator microservice
make run-service

# Run the calculator API client (requires service to be running)
make run-client

# Run the demo (starts both service and client)
./demo.sh
```

See the Makefile for more commands and options:

```bash
make help
```

## Architecture

The project demonstrates a layered architecture:

1. Core logic in packages (`pkg/`)
2. Applications using those packages (`cmd/`)
3. Multiple interfaces to the same functionality (direct CLI, microservice, API client)

This shows different patterns for exposing the same functionality:

- Direct library use
- RESTful API microservice
- API client

## Testing

The project includes examples of:

- Unit tests
- Example-based tests (which also serve as documentation)
- Benchmarks
- Integration tests

### Running Basic Tests

To run unit tests:

```bash
make test
```

### Comprehensive Testing

The project includes a comprehensive test script that:

1. Verifies compilation of all components
2. Runs all unit tests
3. Tests the microservice with both logging systems (ZAP and SLOG)
4. Tests the client-server interaction
5. Performs integration testing across all components
6. Generates detailed code coverage reports

To run the comprehensive tests:

```bash
./test_all.sh
```

### Code Coverage

The core packages have the following test coverage:

| Package | Coverage |
|---------|----------|
| pkg/calculator | 96.6% |
| pkg/logger | 55.2% |
| pkg/slogger | 100.0% |

The test script generates both a console coverage report and an HTML report (`coverage.html`) that provides a visual representation of code coverage.

#### Coverage Analysis

1. **calculator package (96.6%)**: 
   - All primary calculation functions have 100% coverage
   - The no-op logger methods are not covered but are simple pass-through implementations

2. **logger package (55.2%)**:
   - Constructor functions are well-tested
   - Instance methods have minimal coverage since they mostly delegate to the underlying zap logger

3. **slogger package (100%)**:
   - Complete coverage of all functionality

The command-line applications (CLI app, microservice, client) don't have dedicated unit tests as they're primarily tested through integration testing in the `test_all.sh` script.

## License

MIT