# Makefile for iPhone Service API

# Variables
APP_NAME=iphone-service-api
# Docker commands - can be overridden with DOCKER_SUDO=sudo if needed
# Note: To avoid using sudo, add your user to docker group: sudo usermod -aG docker $USER
DOCKER_SUDO ?= 
# Add space after DOCKER_SUDO if it's not empty
ifeq ($(DOCKER_SUDO),)
  DOCKER=docker
  DOCKER_COMPOSE_CMD := $(shell command -v docker-compose >/dev/null 2>&1 && echo "docker-compose" || echo "docker compose")
else
  DOCKER=$(DOCKER_SUDO) docker
  DOCKER_COMPOSE_CMD := $(shell command -v docker-compose >/dev/null 2>&1 && echo "$(DOCKER_SUDO) docker-compose" || echo "$(DOCKER_SUDO) docker compose")
endif
DOCKER_COMPOSE=$(DOCKER_COMPOSE_CMD)

# Go compiler detection - can be overridden with GO=path/to/go
# Prefer standard Go compiler over gccgo for better generics support
ifeq ($(GO),)
  # Try to find standard Go compiler in common locations
  ifeq ($(shell test -x /usr/local/go/bin/go && echo yes),yes)
    GO := /usr/local/go/bin/go
  else ifeq ($(shell test -x $(HOME)/sdk/go1.21/bin/go && echo yes),yes)
    GO := $(HOME)/sdk/go1.21/bin/go
  else ifeq ($(shell test -x $(HOME)/sdk/go1.20/bin/go && echo yes),yes)
    GO := $(HOME)/sdk/go1.20/bin/go
  else ifeq ($(shell test -x $(HOME)/sdk/go1.19/bin/go && echo yes),yes)
    GO := $(HOME)/sdk/go1.19/bin/go
  else
    # Fall back to system go, but check if it's gccgo
    SYSTEM_GO := $(shell which go 2>/dev/null)
    ifneq ($(SYSTEM_GO),)
      GO := $(SYSTEM_GO)
      GO_IS_GCCGO := $(shell $(SYSTEM_GO) version 2>/dev/null | grep -q gccgo && echo yes || echo no)
      ifeq ($(GO_IS_GCCGO),yes)
        $(warning $(YELLOW)Warning: gccgo detected. Standard Go compiler recommended for generics support.$(NC))
        $(warning $(YELLOW)Install from: https://go.dev/dl/$(NC))
      endif
    else
      GO := go
      $(warning $(YELLOW)Warning: Could not find Go compiler. Please install Go from https://go.dev/dl/$(NC))
    endif
  endif
endif

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-up docker-down migrate seed check-go docker-setup docker-help

# Default target
help: ## Show this help message
	@echo "$(BLUE)iPhone Service API - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Check Go compiler
check-go: ## Check Go compiler version and type
	@echo "$(BLUE)Checking Go compiler...$(NC)"
	@echo "Using: $(GO)"
	@$(GO) version 2>/dev/null || (echo "$(RED)Error: Go compiler not found$(NC)" && exit 1)
	@if $(GO) version 2>/dev/null | grep -q gccgo; then \
		echo "$(YELLOW)Warning: gccgo detected. May have issues with generics.$(NC)"; \
		echo "$(YELLOW)Consider installing standard Go: https://go.dev/dl/$(NC)"; \
	else \
		echo "$(GREEN)Standard Go compiler detected.$(NC)"; \
	fi

# Development commands
build: ## Build the application
	@echo "$(BLUE)Building application...$(NC)"
	@$(GO) version > /dev/null 2>&1 || (echo "$(RED)Error: Go compiler not found. Run 'make check-go' for details.$(NC)" && exit 1)
	@if $(GO) version 2>/dev/null | grep -q gccgo; then \
		echo "$(YELLOW)Warning: Using gccgo. Build may fail with generics.$(NC)"; \
	fi
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
# Note: If you get permission errors, either:
#   1. Run: ./scripts/setup-docker.sh (then logout/login)
#   2. Or use: DOCKER_SUDO=sudo make docker-build
docker-setup: ## Setup Docker permissions (add user to docker group)
	@echo "$(BLUE)Setting up Docker permissions...$(NC)"
	@./scripts/setup-docker.sh

docker-help: ## Show Docker help and check setup
	@echo "$(BLUE)Docker Setup Check$(NC)"
	@echo "Docker command: $(DOCKER)"
	@echo "Docker Compose: $(DOCKER_COMPOSE)"
	@echo ""
	@if groups | grep -q docker; then \
		echo "$(GREEN)✅ User is in docker group$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  User is NOT in docker group$(NC)"; \
		echo "$(YELLOW)   Run: make docker-setup (then logout/login)$(NC)"; \
		echo "$(YELLOW)   Or use: DOCKER_SUDO=sudo make docker-up-build$(NC)"; \
	fi
	@echo ""
	@echo "Checking Docker access..."
	@docker --version >/dev/null 2>&1 || (echo "$(RED)❌ Docker command not found$(NC)" && exit 1)
	@docker compose version >/dev/null 2>&1 || docker-compose --version >/dev/null 2>&1 || (echo "$(RED)❌ Docker Compose not found$(NC)" && exit 1)
	@if docker ps >/dev/null 2>&1; then \
		echo "$(GREEN)✅ Can access Docker daemon$(NC)"; \
		echo "$(GREEN)Docker setup looks good!$(NC)"; \
	else \
		echo "$(RED)❌ Cannot access Docker daemon (permission denied)$(NC)"; \
		echo "$(YELLOW)   Run: make docker-setup (then logout/login)$(NC)"; \
		echo "$(YELLOW)   Or use: DOCKER_SUDO=sudo make docker-up-build$(NC)"; \
		exit 1; \
	fi

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	$(DOCKER) build -t $(APP_NAME) .
	@echo "$(GREEN)Docker image built successfully!$(NC)"

docker-rebuild: ## Rebuild Docker image (no cache)
	@echo "$(BLUE)Rebuilding Docker image (no cache)...$(NC)"
	$(DOCKER) build --no-cache -t $(APP_NAME) .
	@echo "$(GREEN)Docker image rebuilt successfully!$(NC)"

docker-up: ## Start all services with Docker Compose
	@echo "$(BLUE)Starting services with Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Services started successfully!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)MinIO Console: http://localhost:9001$(NC)"
	@echo "$(YELLOW)Run 'make docker-logs' to see logs$(NC)"

docker-up-build: ## Build and start all services with Docker Compose
	@echo "$(BLUE)Building and starting services with Docker Compose...$(NC)"
	@if ! $(DOCKER) ps >/dev/null 2>&1; then \
		echo "$(RED)Error: Cannot access Docker daemon$(NC)"; \
		echo "$(YELLOW)Run: make docker-setup (then logout/login)$(NC)"; \
		echo "$(YELLOW)Or use: DOCKER_SUDO=sudo make docker-up-build$(NC)"; \
		exit 1; \
	fi
	$(DOCKER_COMPOSE) up -d --build
	@echo "$(GREEN)Services built and started successfully!$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)MinIO Console: http://localhost:9001$(NC)"
	@echo "$(YELLOW)Run 'make docker-logs' to see logs$(NC)"

docker-down: ## Stop all services
	@echo "$(BLUE)Stopping services...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped successfully!$(NC)"

docker-down-volumes: ## Stop all services and remove volumes
	@echo "$(BLUE)Stopping services and removing volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)Services stopped and volumes removed successfully!$(NC)"

docker-logs: ## Show logs from all services
	$(DOCKER_COMPOSE) logs -f

docker-logs-app: ## Show logs from app service only
	$(DOCKER_COMPOSE) logs -f app

docker-restart: ## Restart all services
	@echo "$(BLUE)Restarting services...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)Services restarted successfully!$(NC)"

docker-ps: ## Show running containers
	$(DOCKER_COMPOSE) ps

docker-exec: ## Execute command in app container (usage: make docker-exec CMD="sh")
	@if [ -z "$(CMD)" ]; then \
		echo "$(RED)Error: Please provide CMD parameter$(NC)"; \
		echo "$(YELLOW)Example: make docker-exec CMD=\"sh\"$(NC)"; \
		exit 1; \
	fi
	$(DOCKER_COMPOSE) exec app $(CMD)

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
	@$(GO) version > /dev/null 2>&1 || (echo "$(RED)Error: Go compiler not found. Run 'make check-go' for details.$(NC)" && exit 1)
	CGO_ENABLED=0 GOOS=linux $(GO) build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/app/main.go
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
