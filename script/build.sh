#!/bin/bash

echo "Building the application..."

# Compile the Go application
go build -o bin/mlvt cmd/server/main.go

echo "Build complete! Executable created in the bin/ directory."
