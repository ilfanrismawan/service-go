# Makefile for iPhone Service API

# Variables
APP_NAME=iphone-service-api
DOCKER_COMPOSE=docker-compose
GO=go
DOCKER=docker

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-up docker-down migrate seed

# Default target
help: ## Show this help message
	@echo "$(BLUE)iPhone Service API - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "$(BLUE)Building application...$(NC)"
	$(GO) build -o bin/$(APP_NAME) cmd/app/main.go
	@echo "$(GREEN)Build completed!$(NC)"

run: ## Run the application locally
	@echo "$(BLUE)Running application...$(NC)"
	$(GO) run cmd/app/main.go

generate-swagger: ## Generate Swagger documentation
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	@if ! command -v swag &> /dev/null; then \
		echo "$(YELLOW)Installing swag...$(NC)"; \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
	@echo "$(GREEN)Swagger documentation generated successfully!$(NC)"
	@echo "$(YELLOW)API Documentation: http://localhost:8080/swagger/index.html$(NC)"
	@echo "$(YELLOW)API Docs: http://localhost:8080/docs$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Database commands
migrate: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	$(GO) run cmd/migrate/main.go

seed: ## Seed database with sample data
	@echo "$(BLUE)Seeding database...$(NC)"
	$(GO) run cmd/seed/main.go

# Docker commands
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	$(DOCKER) build -t $(APP_NAME) .
	@echo "$(GREEN)Docker image built successfully!$(NC)"

docker-up: ## Start all services with Docker Compose
	@echo "$(BLUE)Starting services with Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Services started successfully!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)MinIO Console: http://localhost:9001$(NC)"

docker-down: ## Stop all services
	@echo "$(BLUE)Stopping services...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped successfully!$(NC)"

docker-logs: ## Show logs from all services
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Restart all services
	@echo "$(BLUE)Restarting services...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)Services restarted successfully!$(NC)"

# Development setup
setup: ## Setup development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@if [ ! -f .env ]; then cp env.example .env; echo "$(GREEN)Created .env file from template$(NC)"; fi
	@echo "$(GREEN)Development environment setup completed!$(NC)"
	@echo "$(YELLOW)Please update .env file with your configuration$(NC)"

# Code quality
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	golangci-lint run

format: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	$(GO) vet ./...

# Cleanup
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleanup completed!$(NC)"

clean-docker: ## Clean Docker containers and images
	@echo "$(BLUE)Cleaning Docker artifacts...$(NC)"
	$(DOCKER_COMPOSE) down -v --remove-orphans
	$(DOCKER) system prune -f
	@echo "$(GREEN)Docker cleanup completed!$(NC)"

# API testing
test-api: ## Test API endpoints
	@echo "$(BLUE)Testing API endpoints...$(NC)"
	@echo "$(YELLOW)Health Check:$(NC)"
	curl -s http://localhost:8080/health | jq .
	@echo ""
	@echo "$(YELLOW)API Documentation: http://localhost:8080/swagger/index.html$(NC)"

# Production commands
prod-build: ## Build for production
	@echo "$(BLUE)Building for production...$(NC)"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/app/main.go
	@echo "$(GREEN)Production build completed!$(NC)"

# Quick start
start: setup docker-up ## Quick start (setup + docker-up)
	@echo "$(GREEN)Application is ready!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)Health: http://localhost:8080/health$(NC)"

# Development workflow
dev: setup docker-up ## Start development environment
	@echo "$(GREEN)Development environment is ready!$(NC)"
	@echo "$(YELLOW)Run 'make run' to start the application locally$(NC)"
	@echo "$(YELLOW)Or run 'make docker-logs' to see logs$(NC)"
