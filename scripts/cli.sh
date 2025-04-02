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

# Function to build CRUD Service
build_crud() {
    echo "Building CRUD Service..."
    cd design/crud-api || { echo "Error: Failed to change directory"; exit 1; }
    
    echo "Running go build ./..."
    go build ./... || { echo "Error: Failed to build packages"; exit 1; }
    
    echo "Building crud-service..."
    go build -o crud-service cmd/server/service.go cmd/server/utils.go || { echo "Error: Failed to build crud-service"; exit 1; }
    
    echo "Running tests..."
    go test -v ./... || { echo "Error: Failed to test packages"; exit 1; }
    
    cd ../../
    echo "CRUD Service build completed successfully!"
}

# Function to build Update Service
build_update() {
    echo "Building Update Service..."
    cd design/update-api || { echo "Error: Failed to change directory"; exit 1; }
    
    echo "Running bal test"
    bal test || { echo "Error: Failed to test packages"; exit 1; }
    
    cd ../../
    echo "Update Service build completed successfully!"
}

# Function to build Query Service
build_query() {
    echo "Building Query Service..."
    cd design/query-api || { echo "Error: Failed to change directory"; exit 1; }
    
    echo "Running bal test"
    bal test || { echo "Error: Failed to test packages"; exit 1; }
    
    cd ../../
    echo "Query Service build completed successfully!"
}

# Function to run CRUD Service
run_crud() {
    echo "Starting CRUD Service..."
    cd design/crud-api || { echo "Error: Failed to change directory"; exit 1; }
    ./crud-service &
    CRUD_PID=$!
    PIDS+=("$CRUD_PID")
    echo "CRUD Service started with PID: $CRUD_PID"
    cd ../../
}

# Function to run Update Service
run_update() {
    echo "Starting Update Service..."
    cd design/update-api || { echo "Error: Failed to change directory"; exit 1; }
    bal run . &
    UPDATE_PID=$!
    PIDS+=("$UPDATE_PID")
    echo "Update Service started with PID: $UPDATE_PID"
    cd ../../
}

# Function to run Query Service
run_query() {
    echo "Starting Query Service..."
    cd design/query-api || { echo "Error: Failed to change directory"; exit 1; }
    bal run . &
    QUERY_PID=$!
    PIDS+=("$QUERY_PID")
    echo "Query Service started with PID: $QUERY_PID"
    cd ../../
}

# Function to show help
show_help() {
    echo "Usage: ./cli.sh [command]"
    echo ""
    echo "Commands:"
    echo "  build-crud    Build the CRUD Service"
    echo "  build-update  Build the Update Service"
    echo "  build-query   Build the Query Service"
    echo "  run-crud      Run the CRUD Service"
    echo "  run-update    Run the Update Service"
    echo "  run-query     Run the Query Service"
    echo "  stop          Stop all running services"
    echo "  help          Show this help message"
    echo ""
    echo "Example:"
    echo "  ./cli.sh build-crud"
    echo "  ./cli.sh run-crud"
    echo "  ./cli.sh stop"
}

# Main command handling
case "$1" in
    "build-crud")
        build_crud
        ;;
    "build-update")
        build_update
        ;;
    "build-query")
        build_query
        ;;
    "run-crud")
        run_crud
        ;;
    "run-update")
        run_update
        ;;
    "run-query")
        run_query
        ;;
    "stop")
        stop_services
        ;;
    "help"|"")
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
