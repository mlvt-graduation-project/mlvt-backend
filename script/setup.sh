#!/bin/bash

echo "Setting up the development environment..."

# Install Go dependencies
go mod tidy

# Install other dependencies, e.g., Wire, Mockgen
go get github.com/google/wire/cmd/wire
go get github.com/golang/mock/mockgen

echo "Development environment setup complete!"
