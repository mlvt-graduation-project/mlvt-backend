#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Variables
CMD_DIR="cmd/server"
APP_NAME="mlvt"

# Function to display a message in green
log_info() {
    echo -e "${GREEN}$1${NC}"
}

# Step 1: Install dependencies
log_info "Step 1: Installing Go module dependencies..."
go install github.com/google/wire/cmd/wire@latest
go install github.com/swaggo/swag/cmd/swag@latest
go mod tidy
go mod vendor

# Step 2: Generate Wire dependencies
log_info "Step 2: Generating Wire dependencies..."
cd $CMD_DIR && wire
cd - # Go back to the root directory

# Step 3: Build the application
log_info "Step 3: Building the application..."
cd $CMD_DIR && go build -o $APP_NAME
cd - # Go back to the root directory

# Step 4: Generate Swagger documentation
log_info "Step 4: Generating Swagger documentation..."
swag init -g $CMD_DIR/main.go -o ./docs

# Step 5: Run the built application
log_info "Step 6: Running the built application..."
cd $CMD_DIR && ./$APP_NAME
