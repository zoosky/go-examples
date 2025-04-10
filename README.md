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

### 3. CLI Calculator App

- Located in: `cmd/app`
- Command-line application for basic arithmetic
- Uses the calculator package directly
- Interactive interface

### 4. Calculator Microservice

- Located in: `cmd/calcservice`
- RESTful API for calculator operations
- JSON request/response format
- Configurable port and log level
- Health check endpoint
- Graceful shutdown

### 5. Calculator API Client

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

To run tests:

```bash
make test
```

## License

MIT