.PHONY: help build run test lint fmt deps clean migrate-create migrate-up migrate-down

APP_NAME := backend
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin
PKG := ./...
MIGRATIONS_DIR := ./db/migrations
MIGRATE ?= migrate

# Load local database settings when available. DATABASE_URL takes precedence.
-include .env
DB_SSLMODE ?= $(DB_SSL_MODE)
ifneq ($(strip $(DATABASE_URL)),)
else ifneq ($(strip $(DB_HOST)),)
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
endif

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run $(CMD_DIR)

test: ## Run tests
	@echo "Running tests..."
	@go test -v $(PKG)

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run $(PKG)

fmt: ## Format code
	@echo "Formatting..."
	@go fmt $(PKG)

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

migrate-create: ## Create a new SQL migration (use name=<migration-name>)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=<migration-name>" >&2; exit 1; fi
	@$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up: ## Apply all pending database migrations
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL or database settings in .env are required" >&2; exit 1; fi
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose up

migrate-down: ## Roll back the latest database migration
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL or database settings in .env are required" >&2; exit 1; fi
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" -verbose down 1
