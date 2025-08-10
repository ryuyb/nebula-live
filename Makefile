.PHONY: help build run test clean dev docker-build docker-run docker-dev format lint vet deps tidy check air install-tools swagger-install swagger-gen swagger-validate swagger-serve

# Variables
APP_NAME := nebula-live
BINARY_NAME := server
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/server
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := latest

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
WHITE := \033[37m
RESET := \033[0m

## help: Show this help message
help:
	@echo "$(CYAN)$(APP_NAME) - Makefile Help$(RESET)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(RESET)"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""

## build: Build the application binary
build:
	@echo "$(BLUE)Building $(APP_NAME)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

## run: Run the application
run:
	@echo "$(BLUE)Starting $(APP_NAME)...$(RESET)"
	@go run $(MAIN_PATH)

## dev: Start development server with hot reload (requires Air)
dev:
	@echo "$(BLUE)Starting development server with hot reload...$(RESET)"
	@air -c .air.toml

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

## bench: Run benchmarks
bench:
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	@go test -bench=. ./...

## clean: Clean build artifacts and temporary files
clean:
	@echo "$(BLUE)Cleaning up...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/
	@rm -rf .air/
	@rm -rf data/*.db
	@rm -rf data/*.db-shm
	@rm -rf data/*.db-wal
	@rm -f coverage.out coverage.html
	@rm -f build-errors.log
	@echo "$(GREEN)✓ Cleanup completed$(RESET)"

## format: Format Go code
format:
	@echo "$(BLUE)Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

## lint: Run golangci-lint
lint:
	@echo "$(BLUE)Running linter...$(RESET)"
	@golangci-lint run
	@echo "$(GREEN)✓ Linting completed$(RESET)"

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✓ Vet completed$(RESET)"

## deps: Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(RESET)"

## tidy: Tidy up dependencies
tidy:
	@echo "$(BLUE)Tidying up dependencies...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies tidied$(RESET)"

## check: Run all checks (format, vet, lint, test)
check: format vet lint test
	@echo "$(GREEN)✓ All checks passed$(RESET)"

## install-tools: Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(RESET)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Development tools installed$(RESET)"

## swagger-install: Install Swagger code generation tool
swagger-install:
	@echo "$(BLUE)Installing Swagger tools...$(RESET)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✓ Swagger tools installed$(RESET)"

## swagger-gen: Generate Swagger documentation
swagger-gen:
	@echo "$(BLUE)Generating Swagger documentation...$(RESET)"
	@swag init -g docs.go --output ./docs
	@echo "$(GREEN)✓ Swagger documentation generated in ./docs/$(RESET)"

## swagger-validate: Validate Swagger documentation
swagger-validate:
	@echo "$(BLUE)Validating Swagger documentation...$(RESET)"
	@swag init -g docs.go --output ./docs --parseVendor
	@echo "$(GREEN)✓ Swagger documentation validation completed$(RESET)"

## swagger-serve: Serve Swagger UI locally (requires swagger-ui-dist)
swagger-serve:
	@echo "$(BLUE)Starting local Swagger UI server...$(RESET)"
	@echo "$(YELLOW)Please start the application with 'make run' or 'make dev' first$(RESET)"
	@echo "$(CYAN)Swagger UI available at: http://localhost:8080/swagger/index.html$(RESET)"
	@echo "$(CYAN)Swagger JSON available at: http://localhost:8080/swagger/doc.json$(RESET)"

## docker-build: Build Docker image
docker-build:
	@echo "$(BLUE)Building Docker image...$(RESET)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(RESET)"

## docker-build-dev: Build development Docker image
docker-build-dev:
	@echo "$(BLUE)Building development Docker image...$(RESET)"
	@docker build -f Dockerfile.dev -t $(DOCKER_IMAGE):dev .
	@echo "$(GREEN)✓ Development Docker image built: $(DOCKER_IMAGE):dev$(RESET)"

## docker-run: Run application in Docker container
docker-run:
	@echo "$(BLUE)Running Docker container...$(RESET)"
	@docker run -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-run-dev: Run development Docker container with hot reload
docker-run-dev:
	@echo "$(BLUE)Running development Docker container...$(RESET)"
	@docker-compose -f docker-compose.dev.yml up app-dev

## docker-stop: Stop and remove Docker container
docker-stop:
	@echo "$(BLUE)Stopping Docker container...$(RESET)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true

## compose-up: Start services with Docker Compose
compose-up:
	@echo "$(BLUE)Starting services with Docker Compose...$(RESET)"
	@docker-compose up -d

## compose-up-full: Start full stack (with PostgreSQL and Redis)
compose-up-full:
	@echo "$(BLUE)Starting full stack with Docker Compose...$(RESET)"
	@docker-compose --profile postgres --profile redis up -d

## compose-down: Stop Docker Compose services
compose-down:
	@echo "$(BLUE)Stopping Docker Compose services...$(RESET)"
	@docker-compose down

## compose-logs: Show Docker Compose logs
compose-logs:
	@docker-compose logs -f

## db-sqlite: Switch to SQLite configuration
db-sqlite:
	@echo "$(BLUE)Switching to SQLite configuration...$(RESET)"
	@cp configs/config-sqlite.yaml configs/config.yaml
	@echo "$(GREEN)✓ Switched to SQLite database$(RESET)"

## db-reset: Reset database (remove SQLite file)
db-reset:
	@echo "$(BLUE)Resetting database...$(RESET)"
	@rm -f data/nebula_live.db*
	@echo "$(GREEN)✓ Database reset$(RESET)"

## logs: Show application logs (if running with Docker Compose)
logs:
	@docker-compose logs -f app

## health: Check application health
health:
	@echo "$(BLUE)Checking application health...$(RESET)"
	@curl -s http://localhost:8080/health | jq . || echo "Application not running or jq not installed"

## release: Build release version
release: clean test
	@echo "$(BLUE)Building release version...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✓ Release builds completed$(RESET)"

## info: Show project information
info:
	@echo "$(CYAN)Project Information:$(RESET)"
	@echo "  Name: $(APP_NAME)"
	@echo "  Binary: $(BINARY_NAME)"
	@echo "  Go version: $(shell go version)"
	@echo "  Build dir: $(BUILD_DIR)"
	@echo "  Main path: $(MAIN_PATH)"
	@echo "  Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@echo ""
	@echo "$(YELLOW)Quick Start:$(RESET)"
	@echo "  make db-sqlite    # Switch to SQLite"
	@echo "  make dev          # Start development server"
	@echo "  make health       # Check health"
	@echo "  make swagger-gen  # Generate API docs"
	@echo ""
	@echo "$(YELLOW)Swagger Commands:$(RESET)"
	@echo "  make swagger-install   # Install Swagger tools"
	@echo "  make swagger-gen       # Generate documentation"
	@echo "  make swagger-validate  # Validate documentation"
	@echo "  make swagger-serve     # Show Swagger URLs"