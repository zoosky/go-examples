#!/bin/bash

# Demo script for calculator service and client

echo "Building calculator service and client..."
make build-service
make build-client

# Start the service in the background
echo "Starting calculator service on port 8080..."
./calcservice &
SERVICE_PID=$!

# Give the service time to start
sleep 2

# Run the client
echo "Starting calculator client..."
./calcclient

# When the client exits, terminate the service
echo "Stopping calculator service..."
kill $SERVICE_PID