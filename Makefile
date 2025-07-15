# Go API Makefile

# Variables
BINARY_NAME=api
MAIN_PATH=cmd/main.go
BUILD_DIR=bin

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
.PHONY: run
run: ## Run the application
	go run $(MAIN_PATH)

.PHONY: build
build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: build-linux
build-linux: ## Build for Linux
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)

# Testing commands
.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Development tools
.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	golangci-lint run

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

# Development server
.PHONY: dev
dev: ## Run with hot reload (requires air to be installed)
	air

.PHONY: install-dev-tools
install-dev-tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands (if you want to add Docker support later)
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 4260:4260 $(BINARY_NAME)
