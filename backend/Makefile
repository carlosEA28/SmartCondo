.PHONY: help build run test lint fmt deps clean

APP_NAME := backend
CMD_DIR := ./cmd/api
BUILD_DIR := ./bin
PKG := ./...

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

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
