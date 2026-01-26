.PHONY: help setup-hooks test lint fmt vet build run clean install-tools

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup-hooks: ## Install pre-commit hooks
	@bash scripts/setup-hooks.sh

install-tools: ## Install required development tools
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed successfully!"

test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

lint: ## Run linting checks
	@echo "Running linters..."
	@go fmt ./...
	@go vet ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

build: ## Build the application
	@echo "Building..."
	@go build -o bin/server ./cmd/server

run: ## Run the application
	@echo "Running server..."
	@go run ./cmd/server/main.go

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go -o docs

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out
	@go clean

all: lint test build ## Run lint, test, and build
