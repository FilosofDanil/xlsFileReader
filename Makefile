# Makefile for Go XLS File Reader Bot

# Variables
APP_NAME=xls-reader-bot
DOCKER_IMAGE=xls-reader-bot
DOCKER_TAG=latest
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_TEST=$(GO_CMD) test
GO_CLEAN=$(GO_CMD) clean
GO_MOD=$(GO_CMD) mod
BINARY_NAME=bot
MAIN_PATH=./cmd/main.go

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run clean test docker-build docker-run docker-stop docker-clean deps lint format

# Default target
help: ## Show this help message
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	$(GO_BUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BINARY_NAME)$(NC)"

run: ## Run the application
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	$(GO_RUN) $(MAIN_PATH)

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -rf files/*.xls files/*.xlsx files/debug.txt files/script.txt
	rm -rf log/*.txt logs/*.log
	@echo "$(GREEN)Clean complete$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO_TEST) -v ./...

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO_MOD) download
	$(GO_MOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

format: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)Format complete$(NC)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -d \
		--name $(APP_NAME) \
		--restart unless-stopped \
		-v $(PWD)/files:/app/files \
		-v $(PWD)/logs:/app/logs \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)Container started: $(APP_NAME)$(NC)"
	@echo "$(YELLOW)View logs with: docker logs -f $(APP_NAME)$(NC)"

docker-stop: ## Stop Docker container
	@echo "$(YELLOW)Stopping container...$(NC)"
	docker stop $(APP_NAME)
	docker rm $(APP_NAME)
	@echo "$(GREEN)Container stopped$(NC)"

docker-logs: ## Show Docker container logs
	docker logs -f $(APP_NAME)

docker-clean: ## Remove Docker image
	@echo "$(YELLOW)Removing Docker image...$(NC)"
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)Docker image removed$(NC)"

docker-rebuild: docker-stop docker-clean docker-build docker-run ## Rebuild and restart Docker container

install: deps build ## Install dependencies and build

dev: ## Run in development mode with auto-reload (requires air)
	@echo "$(GREEN)Starting development server...$(NC)"
	air

all: clean deps build ## Clean, install deps and build

.DEFAULT_GOAL := help


