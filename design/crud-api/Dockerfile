# Use the official Go image as the base image
FROM golang:1.24

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o crud-service cmd/server/service.go cmd/server/utils.go

# Expose the port the app runs on
EXPOSE 50051

# Command to run the executable
CMD ["./crud-service"] 