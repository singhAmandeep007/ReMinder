# ============================================
# ReMinder NestJS Server - Makefile
# ============================================
# Usage: make <command>
# Run 'make help' for available commands
# ============================================

.PHONY: help setup install dev start build clean \
        db-up db-down db-reset db-logs \
        lint format typecheck \
        test test-watch test-cov test-e2e \
        docker-build docker-run logs

# Default target
.DEFAULT_GOAL := help

# Colors for terminal output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

# ============================================
# Help
# ============================================
help: ## Show this help message
	@echo ""
	@echo "$(CYAN)ReMinder NestJS Server$(RESET)"
	@echo "$(CYAN)======================$(RESET)"
	@echo ""
	@echo "$(GREEN)Available commands:$(RESET)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-15s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ============================================
# Setup & Installation
# ============================================
setup: install db-up dev ## Full setup: install deps, start DB, run dev server

install: ## Install npm dependencies
	@echo "$(GREEN)Installing dependencies...$(RESET)"
	npm install

# ============================================
# Development
# ============================================
dev: ## Start development server with hot reload
	@echo "$(GREEN)Starting development server...$(RESET)"
	npm run start:dev

start: ## Start server (production mode)
	@echo "$(GREEN)Starting production server...$(RESET)"
	npm run start:prod

build: ## Build for production
	@echo "$(GREEN)Building application...$(RESET)"
	npm run build

clean: ## Clean build artifacts and dependencies
	@echo "$(YELLOW)Cleaning...$(RESET)"
	rm -rf dist node_modules coverage .turbo

# ============================================
# Database (Docker)
# ============================================
db-up: ## Start MongoDB container
	@echo "$(GREEN)Starting MongoDB...$(RESET)"
	docker-compose up -d mongodb
	@echo "$(GREEN)MongoDB is running on port 27017$(RESET)"

db-down: ## Stop MongoDB container
	@echo "$(YELLOW)Stopping MongoDB...$(RESET)"
	docker-compose down

db-reset: db-down ## Reset MongoDB (stop, remove volume, start fresh)
	@echo "$(YELLOW)Removing MongoDB data...$(RESET)"
	docker volume rm nestjs-server_mongodb_data 2>/dev/null || true
	@$(MAKE) db-up

db-logs: ## Show MongoDB logs
	docker-compose logs -f mongodb

db-shell: ## Open MongoDB shell
	@echo "$(GREEN)Opening MongoDB shell...$(RESET)"
	docker exec -it reminder-mongodb mongosh

# ============================================
# Code Quality
# ============================================
lint: ## Run ESLint
	@echo "$(GREEN)Running linter...$(RESET)"
	npm run lint

format: ## Format code with Prettier
	@echo "$(GREEN)Formatting code...$(RESET)"
	npm run format

typecheck: ## Run TypeScript type check
	@echo "$(GREEN)Running type check...$(RESET)"
	npm run typecheck

check: lint typecheck ## Run all code quality checks

# ============================================
# Testing
# ============================================
test: ## Run unit tests
	@echo "$(GREEN)Running unit tests...$(RESET)"
	npm run test

test-watch: ## Run unit tests in watch mode
	@echo "$(GREEN)Running tests in watch mode...$(RESET)"
	npm run test:watch

test-cov: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(RESET)"
	npm run test:cov

test-e2e: ## Run E2E tests
	@echo "$(GREEN)Running E2E tests...$(RESET)"
	npm run test:e2e

test-all: test test-e2e ## Run all tests

# ============================================
# Type Generation
# ============================================
generate-types: ## Generate TypeScript types from OpenAPI spec
	@echo "$(GREEN)Generating types from OpenAPI spec...$(RESET)"
	npm run generate:types

# ============================================
# Docker
# ============================================
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(RESET)"
	docker build -t reminder-api .

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(RESET)"
	docker run -p 3000:3000 \
		-e MONGO_DB_URI=mongodb://host.docker.internal:27017 \
		-e MONGO_DB_NAME=reminder \
		reminder-api

# ============================================
# Utilities
# ============================================
logs: ## Tail application logs
	@tail -f logs/*.log 2>/dev/null || echo "No logs found. Start the server first."

kill-port: ## Kill process on port 3000
	@echo "$(YELLOW)Killing process on port 3000...$(RESET)"
	@lsof -ti:3000 | xargs kill -9 2>/dev/null || echo "No process on port 3000"

health: ## Check server health
	@curl -s http://localhost:3000/api/health/live | jq . 2>/dev/null || \
		echo "Server not running or jq not installed"

# ============================================
# Composite Commands
# ============================================
fresh: clean install db-reset ## Clean install with fresh database

ci: check test build ## Run CI pipeline locally

all: setup ## Alias for setup
