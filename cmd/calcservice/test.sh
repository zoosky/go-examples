#!/bin/bash

# Simple test script for calculator microservice
echo "Testing calculator microservice..."

# Start the service in the background
./calcservice &
PID=$!

# Wait for the service to start
sleep 1
echo "Service started with PID $PID"

# Test health endpoint
echo -e "\nTesting health endpoint:"
curl -s http://localhost:8080/health | jq

# Test addition
echo -e "\nTesting addition (5 + 3):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 5, "b": 3}' | jq

# Test subtraction
echo -e "\nTesting subtraction (10 - 4):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "subtract", "a": 10, "b": 4}' | jq

# Test multiplication
echo -e "\nTesting multiplication (6 * 7):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "multiply", "a": 6, "b": 7}' | jq

# Test division
echo -e "\nTesting division (20 / 5):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 20, "b": 5}' | jq

# Test division by zero
echo -e "\nTesting division by zero (10 / 0):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 10, "b": 0}' | jq

# Test invalid operation
echo -e "\nTesting invalid operation (power):"
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "power", "a": 2, "b": 3}' | jq

# Kill the service
echo -e "\nStopping service..."
kill $PID

echo "Test completed"