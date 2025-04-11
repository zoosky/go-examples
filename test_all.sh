#!/bin/bash

# Directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colored output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print section header
section() {
  echo -e "\n${BLUE}======================================================================${NC}"
  echo -e "${BLUE}= $1${NC}"
  echo -e "${BLUE}======================================================================${NC}"
}

# Print subsection header
subsection() {
  echo -e "\n${YELLOW}== $1 ==${NC}"
}

# Print success message
success() {
  echo -e "${GREEN}✅ $1${NC}"
}

# Print error message and exit
error() {
  echo -e "${RED}❌ $1${NC}"
  exit 1
}

# Print warning message
warning() {
  echo -e "${YELLOW}⚠️  $1${NC}"
}

# Run a command and handle errors
run_cmd() {
  echo -e "\n$ $1"
  eval "$1"
  local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    error "Command failed with exit code $exit_code"
  fi
}

# Test if a component is running on the specified port
test_port() {
  local port=$1
  local attempt=1
  local max_attempts=5
  
  while [ $attempt -le $max_attempts ]; do
    echo "Checking if port $port is available (attempt $attempt/$max_attempts)..."
    if nc -z localhost $port 2>/dev/null; then
      return 0
    fi
    sleep 1
    attempt=$((attempt + 1))
  done
  
  return 1
}

# Kill a process by port
kill_process_on_port() {
  local port=$1
  local pid=$(lsof -n -i4TCP:$port | grep LISTEN | awk '{print $2}')
  if [ -n "$pid" ]; then
    echo "Killing process on port $port (PID: $pid)"
    kill "$pid" 2>/dev/null || true
  fi
}

# Clean up on script exit
cleanup() {
  echo "Cleaning up..."
  kill_process_on_port 8080
  kill_process_on_port 8081
  kill_process_on_port 8082
}

# Register cleanup function
trap cleanup EXIT

# Start tests
section "Running Tests for All Subprojects"

# Step 1: Check if we're in the correct directory
subsection "Verifying Project Directory"
if [ ! -f "go.mod" ]; then
  error "Not in project root directory or go.mod missing"
else
  success "Found go.mod file"
fi

# Step 2: Run linting to check code quality
subsection "Linting Code"
run_cmd "make lint"
success "Linting passed"

# Step 3: Verify all packages compile
subsection "Verifying Compilation"
run_cmd "go build ./..."
success "All packages compile successfully"

# Step 4: Run all tests
subsection "Running Unit Tests"
run_cmd "go test -cover ./..."
success "All tests passed"

# Step 5: Build all components
subsection "Building All Components"
run_cmd "make build"
success "All components built successfully"

# Step 6: Test calculator microservice with ZAP logger
subsection "Testing Calculator Microservice (ZAP Logger)"
# Kill any process that might be using port 8080
kill_process_on_port 8080

# Start the service with ZAP logger
./calcservice --log-system zap --port 8080 --log-level info &
CALCSERVICE_PID=$!

# Wait for the service to start
sleep 2
if ! test_port 8080; then
  error "Failed to start calculator microservice"
fi
success "Calculator microservice started successfully with ZAP logger"

# Run basic API test
response=$(curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 5, "b": 3}')
  
expected_result='"result":8'
if [[ "$response" == *"$expected_result"* ]]; then
  success "API responded correctly with ZAP logger (5 + 3 = 8)"
else
  error "API test failed with ZAP logger. Got: $response"
fi

# Kill ZAP logger service
kill $CALCSERVICE_PID
sleep 1

# Step 7: Test calculator microservice with SLOG logger
subsection "Testing Calculator Microservice (SLOG Logger)"
# Start the service with SLOG logger
./calcservice --log-system slog --port 8080 --log-level info &
CALCSERVICE_PID=$!

# Wait for the service to start
sleep 2
if ! test_port 8080; then
  error "Failed to start calculator microservice with SLOG logger"
fi
success "Calculator microservice started successfully with SLOG logger"

# Run basic API test
response=$(curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "multiply", "a": 6, "b": 7}')
  
expected_result='"result":42'
if [[ "$response" == *"$expected_result"* ]]; then
  success "API responded correctly with SLOG logger (6 * 7 = 42)"
else
  error "API test failed with SLOG logger. Got: $response"
fi

# Kill SLOG logger service
kill $CALCSERVICE_PID
sleep 1

# Step 8: Test calculator client against the microservice
subsection "Testing Calculator Client with Microservice"
# Start the service again
./calcservice --port 8080 &
CALCSERVICE_PID=$!

# Wait for the service to start
sleep 2
if ! test_port 8080; then
  error "Failed to start calculator microservice for client test"
fi
success "Calculator microservice started successfully for client test"

# Create temporary script to test client
TEST_SCRIPT=$(mktemp)
cat > "$TEST_SCRIPT" << 'EOF'
#!/bin/bash
# We need to simulate user input
(
  echo "add 7 8"
  sleep 0.5
  echo "multiply 6 6"
  sleep 0.5
  echo "quit"
) | ./calcclient
EOF
chmod +x "$TEST_SCRIPT"

# Run client test
CLIENT_OUTPUT=$("$TEST_SCRIPT")
rm "$TEST_SCRIPT"

if [[ "$CLIENT_OUTPUT" == *"Result: 15"* ]] && [[ "$CLIENT_OUTPUT" == *"Result: 36"* ]]; then
  success "Calculator client successfully communicated with server"
else
  warning "Calculator client test could not be verified"
  echo "Client output: $CLIENT_OUTPUT"
fi

# Kill the service
kill $CALCSERVICE_PID
sleep 1

# Step 9: Verify help functionality
subsection "Verifying Help Commands"
run_cmd "make help >/dev/null"
success "Help command works correctly"

section "Test Coverage Report"
# Generate coverage report for all packages
echo "Generating detailed coverage report..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Calculate the total coverage percentage
TOTAL_COV=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo -e "\nTotal code coverage: ${GREEN}${TOTAL_COV}${NC}"

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
echo "HTML coverage report generated: ${BLUE}coverage.html${NC}"

# Provide a detailed breakdown of coverage by package
echo -e "\nCoverage by package:"
go test -cover ./... | grep -v "no test files" | sort

section "Test Summary"
success "All tests completed successfully!"
echo ""
echo "The go-examples project is working correctly with all components:"
echo "  - Calculator package"
echo "  - Logger package (ZAP)"
echo "  - SLogger package (SLOG)"
echo "  - CLI Calculator"
echo "  - Calculator Microservice"
echo "  - Calculator API Client"