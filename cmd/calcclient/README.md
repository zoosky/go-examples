# Calculator Client

This is a command-line client for the calculator microservice. It uses the microservice API to perform calculator operations.

## Features

- Command-line interface for calculator operations
- Connects to the calculator microservice
- Support for add, subtract, multiply, and divide operations
- Connection health check
- Configurable server URL and timeout

## Usage

### Build and Run

You can build and run the client using Make:

```bash
# Build the client
make build-client

# Run with default settings (connecting to http://localhost:8080)
make run-client

# Run with custom settings
make run-client-custom SERVER=http://example.com:9000 TIMEOUT=10
```

Or build and run manually:

```bash
go build -o calcclient ./cmd/calcclient
./calcclient --server http://localhost:8080 --timeout 5
```

### Command Line Arguments

- `--server`: The URL of the calculator service (default: "http://localhost:8080")
- `--timeout`: Request timeout in seconds (default: 5)

### Interactive Commands

Once the client is running, you can use the following interactive commands:

- `add <number1> <number2>`: Add two numbers
- `subtract <number1> <number2>`: Subtract the second number from the first
- `multiply <number1> <number2>`: Multiply two numbers
- `divide <number1> <number2>`: Divide the first number by the second
- `quit`, `exit`, or `q`: Exit the client

## Examples

```
Calculator Client
================
Connected to: http://localhost:8080
Available operations: add, subtract, multiply, divide, quit
Example usage: add 5 3

> add 5 3
Executing: add 5 3
Result: 8

> multiply 6 7
Executing: multiply 6 7
Result: 42

> divide 10 0
Executing: divide 10 0
Error: API error: Division by zero

> quit
Goodbye!
```

## Notes

- The client requires the calculator microservice to be running
- The client automatically checks if the service is available on startup
- For best performance, run the service and client on the same machine