# Use the official Golang image as the base image
FROM golang:1.20-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application (binary) for the server
RUN go build -o main cmd/server/main.go

# Use a minimal image to run the compiled Go binary
FROM alpine:3.18

# Set the working directory in the new container
WORKDIR /app

# Copy the binary and other necessary files from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/mlvt.db .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]

