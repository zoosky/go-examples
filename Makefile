# Go parameters
BINARY_NAME=app
SERVICE_NAME=calcservice
CLIENT_NAME=calcclient
MAIN_PACKAGE=./cmd/app
SERVICE_PACKAGE=./cmd/calcservice
CLIENT_PACKAGE=./cmd/calcclient
COVERAGE_PROFILE=coverage.out
BUILD_FLAGS=-ldflags="-s -w" -trimpath
GOBIN=$(CURDIR)/bin

# Default make command
.PHONY: all
all: clean lint test build

#################################################
# Build commands
#################################################

# Build all applications
.PHONY: build
build: build-app build-service build-client

# Build the CLI application
.PHONY: build-app
build-app:
	@echo "Building $(BINARY_NAME)..."
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build completed: $(BINARY_NAME)"

# Build the calculator microservice
.PHONY: build-service
build-service:
	@echo "Building $(SERVICE_NAME)..."
	@go build $(BUILD_FLAGS) -o $(SERVICE_NAME) $(SERVICE_PACKAGE)
	@echo "Build completed: $(SERVICE_NAME)"

# Build the calculator client
.PHONY: build-client
build-client:
	@echo "Building $(CLIENT_NAME)..."
	@go build $(BUILD_FLAGS) -o $(CLIENT_NAME) $(CLIENT_PACKAGE)
	@echo "Build completed: $(CLIENT_NAME)"

# Build for multiple platforms (cross-compilation)
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Multi-platform build completed"

# Install binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing binary..."
	@go install $(BUILD_FLAGS) $(MAIN_PACKAGE)
	@echo "Installation completed"

#################################################
# Test commands
#################################################

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "Tests completed"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...
	@echo "Race detection tests completed"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=$(COVERAGE_PROFILE) ./...
	@go tool cover -html=$(COVERAGE_PROFILE)
	@echo "Coverage tests completed"

# Run benchmark tests
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...
	@echo "Benchmarks completed"

#################################################
# Quality control
#################################################

# Format code using gofmt
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .
	@echo "Formatting completed"

# Lint code using golangci-lint
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.56.2; \
		$(GOBIN)/golangci-lint run ./...; \
	fi
	@echo "Linting completed"

# Verify code with multiple checks
.PHONY: verify
verify: fmt lint test

#################################################
# Dependency management
#################################################

# Update dependencies
.PHONY: deps
deps:
	@echo "Updating dependencies..."
	@go mod tidy
	@echo "Dependencies updated"

# Install necessary tools
.PHONY: tools
tools:
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installation completed"

#################################################
# Clean up
#################################################

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-linux-amd64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-windows-amd64.exe
	@rm -f $(SERVICE_NAME) $(SERVICE_NAME)-linux-amd64 $(SERVICE_NAME)-darwin-amd64 $(SERVICE_NAME)-windows-amd64.exe
	@rm -f $(CLIENT_NAME) $(CLIENT_NAME)-linux-amd64 $(CLIENT_NAME)-darwin-amd64 $(CLIENT_NAME)-windows-amd64.exe
	@rm -f $(COVERAGE_PROFILE)
	@echo "Clean completed"

#################################################
# Code generation
#################################################

# Generate code (mocks, swagger, etc)
.PHONY: generate
generate:
	@echo "Generating code..."
	@go generate ./...
	@echo "Code generation completed"

# Run the CLI application
.PHONY: run
run: build-app
	@echo "Running CLI application..."
	@./$(BINARY_NAME)

# Run the calculator microservice
.PHONY: run-service
run-service: build-service
	@echo "Running calculator microservice on port 8080..."
	@./$(SERVICE_NAME)

# Run the service with custom port
.PHONY: run-service-custom
run-service-custom: build-service
	@echo "Running calculator microservice with custom settings..."
	@./$(SERVICE_NAME) --port $(PORT) --log-level $(LOG_LEVEL)

# Run the calculator client (requires the service to be running)
.PHONY: run-client
run-client: build-client
	@echo "Running calculator client (connecting to localhost:8080)..."
	@./$(CLIENT_NAME)

# Run the client with custom server URL
.PHONY: run-client-custom
run-client-custom: build-client
	@echo "Running calculator client with custom settings..."
	@./$(CLIENT_NAME) --server $(SERVER) --timeout $(TIMEOUT)

# Help command to display available make targets
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make              : Build, test, and verify code"
	@echo "  make build        : Build all applications"
	@echo "  make build-app    : Build only the CLI application"
	@echo "  make build-service: Build only the microservice"
	@echo "  make build-client : Build only the client application"
	@echo "  make build-all    : Build for multiple platforms"
	@echo "  make install      : Install binary to GOPATH/bin"
	@echo "  make test         : Run all tests"
	@echo "  make test-race    : Run tests with race detection"
	@echo "  make test-coverage: Run tests with coverage report"
	@echo "  make benchmark    : Run benchmark tests"
	@echo "  make fmt          : Format code using gofmt"
	@echo "  make lint         : Lint code using golangci-lint"
	@echo "  make verify       : Run fmt, lint, and tests"
	@echo "  make deps         : Update dependencies"
	@echo "  make tools        : Install development tools"
	@echo "  make clean        : Clean build artifacts"
	@echo "  make generate     : Generate code (mocks, swagger, etc)"
	@echo "  make run          : Build and run the CLI application"
	@echo "  make run-service  : Build and run the calculator microservice"
	@echo "  make run-service-custom PORT=8080 LOG_LEVEL=debug : Run service with custom settings"
	@echo "  make run-client   : Build and run the calculator client (requires service to be running)"
	@echo "  make run-client-custom SERVER=http://localhost:8080 TIMEOUT=5 : Run client with custom settings"
	@echo "  make help         : Display this help message"

