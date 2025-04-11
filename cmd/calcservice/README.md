# Calculator Microservice

This is a simple HTTP microservice that provides calculator functionality through a REST API. The service uses the calculator package to perform operations, with proper logging.

## Features

- RESTful API for calculator operations
- Support for add, subtract, multiply, and divide operations
- Health check endpoint
- Configurable port and log level
- Multiple logging system options (zap or slog)
- Graceful shutdown on interrupt signal

## Usage

### Build and Run

You can build and run the microservice using Make:

```bash
# Build the service
make build-service

# Run with default settings (port 8080, info log level, zap logger)
make run-service

# Run with custom settings
make run-service-custom PORT=9000 LOG_LEVEL=debug LOG_SYSTEM=slog
```

Or build and run manually:

```bash
go build -o calcservice ./cmd/calcservice
./calcservice --port 8080 --log-level info --log-system zap
```

### Logging Systems

The service supports two logging systems:

1. **zap** (default): Uses Uber's zap logging library for structured, high-performance logging
2. **slog**: Uses Go's standard library slog package for structured logging

To switch between them, use the `--log-system` flag:

```bash
# Use zap logger
./calcservice --log-system zap

# Use slog logger
./calcservice --log-system slog
```

### API Endpoints

#### Calculate

Perform a calculation operation.

- **URL**: `/calculate`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "operation": "add",  // One of: add, subtract, multiply, divide
    "a": 10,
    "b": 5
  }
  ```
- **Success Response**:
  ```json
  {
    "result": 15,
    "success": true
  }
  ```
- **Error Response**:
  ```json
  {
    "success": false,
    "error": "Division by zero"
  }
  ```

#### Health Check

Check if the service is running.

- **URL**: `/health`
- **Method**: `GET`
- **Success Response**:
  ```json
  {
    "status": true
  }
  ```

## Examples

### Using curl

```bash
# Addition
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 5, "b": 3}'

# Subtraction
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "subtract", "a": 10, "b": 4}'

# Multiplication
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "multiply", "a": 6, "b": 7}'

# Division
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 20, "b": 5}'

# Health check
curl http://localhost:8080/health
```