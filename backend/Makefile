.PHONY: help build run test lint migrate migrate-up migrate-down clean

APP_NAME := backend
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run $(CMD_DIR)/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./internal/tests/...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@migrate create -ext sql -dir db/migrations -seq $(name)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

fmt: ## Format code
	@echo "Formatting..."
	@gofmt -s -w .

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
