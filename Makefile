# Load .env (for DATABASE_URL)
include .env
export


# ─── Configuration ─────────────────────────────────────────────────────────────
CMD_DIR := cmd/server
SEED_GO_DIR := cmd/seeder
SQL_SEED_DIR := ./seed
CLEAN_DIR := cmd/cleanup
APP_NAME := mlvt
SCRIPT_DIR := script/
MIGRATION_DIR := migration/Postgres
INITIALIZE_DIR := internal/initialize
OUTPUT_DIR := internal/wire_gen
DATABASE_URL := $(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=${SSL_MODE}


# ─── Default ───────────────────────────────────────────────────────────────────
all: build

# ─── Run Targets ───────────────────────────────────────────────────────────────
run:
	cd $(CMD_DIR) && go run .

seed:
	cd $(SEED_GO_DIR) && go run .



cleaner:
	cd $(CLEAN_DIR) && go run .

# ─── Migration ─────────────────────────────────────────────────────────────────
migrate-up:
	@echo "🔼 Running migrations up..."
	@migrate -path ${MIGRATION_DIR} -database "$(DATABASE_URL)" up

migrate-down:
	@echo "🔽 Rolling back last migration..."
	@migrate -path ${MIGRATION_DIR} -database "$(DATABASE_URL)" down 1

create-migration:
	@read -p "Enter migration name (snake_case): " name; \
	migrate create -ext sql -dir ${MIGRATION_DIR} $$name

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path ${MIGRATION_DIR} -database "$(DATABASE_URL)" force $$version

seed-postgres:
	@echo "🌱 Seeding SQL files from $(SQL_SEED_DIR)..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		for f in $(SQL_SEED_DIR)/*.sql; do \
			echo "📄 Seeding $$f..."; \
			psql -U $(USER) -d $(DB_NAME) -h $(DB_HOST) -p $(DB_PORT) -f "$$f"; \
		done \
	elif [ "$$(uname)" = "Linux" ]; then \
		for f in $(SQL_SEED_DIR)/*.sql; do \
			echo "📄 Seeding $$f..."; \
			sudo -u postgres psql -d $(DB_NAME) -f "$$f"; \
		done \
	else \
		echo "❌ Unsupported OS for sql-seed"; \
	fi
	@echo "✅ Done seeding SQL."

# ─── Build ─────────────────────────────────────────────────────────────────────
build:
	cd $(CMD_DIR) && go build -o $(APP_NAME)

clean:
	rm -f $(CMD_DIR)/$(APP_NAME)

wire:
	cd $(INITIALIZE_DIR) && go run github.com/google/wire/cmd/wire@latest

wire-build:
	wire build

create-db:
	@echo "📦 Creating user and database '$(DB_NAME)' owned by '$(DB_USER)'..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "🖥️ Detected macOS (Homebrew). Using local user '$(USER)'..."; \
		psql -U $(USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -tc "SELECT 1 FROM pg_roles WHERE rolname='$(DB_USER)'" | grep -q 1 || \
		psql -U $(USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -c "CREATE ROLE $(DB_USER) WITH LOGIN PASSWORD '$(DB_PASSWORD)';"; \
		psql -U $(USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -tc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 || \
		psql -U $(USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "🐧 Detected Linux. Using 'postgres' superuser..."; \
		sudo -u postgres psql -d postgres -tc "SELECT 1 FROM pg_roles WHERE rolname='$(DB_USER)'" | grep -q 1 || \
		sudo -u postgres psql -d postgres -c "CREATE ROLE $(DB_USER) WITH LOGIN PASSWORD '$(DB_PASSWORD)';"; \
		sudo -u postgres psql -d postgres -tc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 || \
		sudo -u postgres psql -d postgres -c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"; \
	else \
		echo "❌ Unsupported OS. Please create the database manually."; \
	fi
	@echo "✅ Done creating database '$(DB_NAME)' for user '$(DB_USER)'"

drop-db:
	@echo "💣 Dropping database '$(DB_NAME)' if exists..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		psql -U $(USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS $(DB_NAME);" || true; \
	elif [ "$$(uname)" = "Linux" ]; then \
		sudo -u postgres psql -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);" || true; \
	else \
		echo "❌ Unsupported OS for drop-db"; \
	fi
	@echo "✅ Dropped database '$(DB_NAME)'"

install-postgres:
	@echo "🛠️ Installing PostgreSQL..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "📦 Detected macOS. Installing via Homebrew..."; \
		brew install postgresql@14; \
		echo "👉 To start PostgreSQL: brew services start postgresql@14"; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "📦 Detected Linux. Installing via apt..."; \
		sudo apt update && sudo apt install -y postgresql postgresql-contrib; \
		echo "👉 To start PostgreSQL: sudo service postgresql start"; \
	else \
		echo "❌ Unsupported OS. Please install PostgreSQL manually."; \
	fi

# ─── Deployment ────────────────────────────────────────────────────────────────
deploy:
	@echo "🚀 Building project..."
	cd $(CMD_DIR) && go build -o $(APP_NAME)
	@echo "🔄 Restarting $(APP_NAME) service..."
	sudo systemctl restart $(APP_NAME)
	@echo "✅ Done."

# ─── Migration ─────────────────────────────────────────────────────────────────
migrate-sqlite:
	migrate create -ext sql -dir migrations create_users_table
	migrate create -ext sql -dir migrations create_videos_table
	migrate create -ext sql -dir migrations create_transcriptions_table
	migrate create -ext sql -dir migrations create_transaction_logs_table
	migrate create -ext sql -dir migrations create_frames_table
	migrate create -ext sql -dir migrations create_audios_table

# ─── Swagger ───────────────────────────────────────────────────────────────────
swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g $(CMD_DIR)/main.go -o ./docs

# ─── Script Group ──────────────────────────────────────────────────────────────
run-all:
	bash $(SCRIPT_DIR)/run_all.sh

# ─── Help ──────────────────────────────────────────────────────────────────────
help:
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "  make run              Run the application"
	@echo "  make seed             Run the Go-based seeder"
	@echo "  make seed-postgres    Seed postgres datbase from .sql files in $(SQL_SEED_DIR)" 
	@echo "  make build            Build the application"
	@echo "  make clean            Clean the generated binaries"
	@echo "  make wire             Generate dependencies with Wire"
	@echo "  make wire-build       Generate Wire dependencies and build the application"
	@echo "  make cleaner          Run cleanup script (Go-based)"
	@echo "  make run-all          Run all scripts in $(SCRIPT_DIR)/run_all.sh"
	@echo "  make swag             Generate Swagger docs from code"
	@echo "  make deploy           Build and restart $(APP_NAME) systemd service"
	@echo "  make create-db        Create DB using 'postgres' user (Linux/Ubuntu)"
	@echo "  make create-db-macos  Create DB using current macOS user (Homebrew)"
	@echo "  make drop-db          Drop the target database if exists"
	@echo "  make help             Show this help message"
