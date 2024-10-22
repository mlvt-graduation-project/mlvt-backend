# Variables
CMD_DIR := cmd/server
OUTPUT_DIR := internal/wire_gen
SEED_DIR := cmd/seeder
CLEAN_DIR := cmd/cleanup
APP_NAME := mlvt
SCRIPT_DIR :=script/
INITIALIZE_DIR := internal/initialize

# Default target
all: build

# Run the application
run:
	cd $(CMD_DIR) && go run .

# Run the seeder
seed:
	cd $(SEED_DIR) && go run .

# Run the cleaner
cleaner:
	cd $(CLEAN_DIR) && go run .

# Build the application
build:
	cd $(CMD_DIR) && go build -o $(APP_NAME)

# Generate wire dependencies
wire:
	go install github.com/google/wire/cmd/wire@latest
	cd $(INITIALIZE_DIR) && wire

# Clean the generated binaries
clean:
	rm -f $(CMD_DIR)/$(APP_NAME)

# Generate wire and then build
wire-build:
	wire build

#migration
migrate:
	migrate create -ext sql -dir migrations create_users_table
	migrate create -ext sql -dir migrations create_videos_table
	migrate create -ext sql -dir migrations create_transcriptions_table
	migrate create -ext sql -dir migrations create_transaction_logs_table
	migrate create -ext sql -dir migrations create_frames_table
	migrate create -ext sql -dir migrations create_audios_table


# Swagger
swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(CMD_DIR)/main.go -o ./docs

# Run all steps from the script
run-all:
	bash $(SCRIPT_DIR)/run_all.sh

# Help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make run         Run the application"
	@echo "  make swag		  Run the swagger"
	@echo "  make build       Build the application"
	@echo "  make wire        Generate dependencies with Wire"
	@echo "  make clean       Clean the generated binaries"
	@echo "  make wire-build  Generate Wire dependencies and build the application"
	@echo "  make help        Show this help message"
