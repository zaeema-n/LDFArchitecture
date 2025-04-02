#!/bin/bash

# Initialize array to store background process PIDs
PIDS=()

# Function to stop all services
stop_services() {
    echo "Stopping all services..."
    for pid in "${PIDS[@]}"; do
        if kill -0 $pid 2>/dev/null; then
            echo "Stopping process $pid..."
            kill -9 $pid 2>/dev/null || true
            wait $pid 2>/dev/null || true
        fi
    done
    # Additional cleanup for any remaining processes
    pkill -f crud-service || true
    pkill -f "bal run" || true
    echo "All services stopped"
}

# build CRUD Service
echo "Building CRUD Service..."
cd design/crud-api && echo "Changed directory to: $(pwd)"

echo "Running go build ./..."
go build ./... || { echo "Error: Failed to build packages"; exit 1; }

echo "Building crud-service..."
go build -o crud-service cmd/server/service.go cmd/server/utils.go || { echo "Error: Failed to build crud-service"; exit 1; }

echo "Build completed successfully!"

echo "Running tests..."
go test -v ./... || { echo "Error: Failed to test packages"; exit 1; }

echo "Tests completed successfully!"

echo "Starting CRUD Service..."
./crud-service &
CRUD_PID=$!
PIDS+=("$CRUD_PID")
echo "CRUD Service started with PID: $CRUD_PID"

cd ../../

# build Update Service
echo "Building Update Service..."
cd design/update-api && echo "Changed directory to: $(pwd)"

echo "Running bal test"
bal test || { echo "Error: Failed to test packages"; exit 1; }

echo "Tests completed successfully!"

echo "Build completed successfully!"

echo "Starting Update Service..."
bal run . &
UPDATE_PID=$!
PIDS+=("$UPDATE_PID")
echo "Update Service started with PID: $UPDATE_PID"

cd ../../

# building Query Service
echo "Building Query Service..."
cd design/query-api && echo "Changed directory to: $(pwd)"

echo "Running bal test"
bal test || { echo "Error: Failed to test packages"; exit 1; }

echo "Tests completed successfully!"

echo "Build completed successfully!"

echo "Starting Query Service..."
bal run . &
QUERY_PID=$!
PIDS+=("$QUERY_PID")
echo "Query Service started with PID: $QUERY_PID"

# Wait for a specified time (e.g., 30 seconds)
echo "Services are running. Waiting for 30 seconds..."
sleep 5

# Stop all services
stop_services

echo "Script completed successfully!"
