.PHONY: help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build web and api binary
	@echo "Building web frontend..."
	@cd web && bun run build
	@echo "Building API binary..."
	@go build -o bin/api ./backend/

dev: ## builds the web, then the api and runs the binary directly and spins up a postgres container
	@echo "Starting development environment..."
	@echo "1. Starting PostgreSQL container..."
	@docker compose -f docker-compose.dev.yaml up -d postgres
	@echo "2. Waiting for PostgreSQL to be ready..."
	@until docker compose -f docker-compose.dev.yaml exec postgres pg_isready -U openwaitlist -d openwaitlist; do \
		echo "Waiting for postgres..."; \
		sleep 2; \
	done
	@echo "3. Building web frontend..."
	@cd web && bun run build
	@echo "4. Building API binary..."
	@go build -o bin/api ./backend/
	@echo "5. Starting API server..."
	@set -a && [ -f ./.env ] && . ./.env && set +a && ./bin/api

dev-stop: ## stops the postgres container
	@echo "Stopping development environment..."
	@docker compose -f docker-compose.dev.yaml down


# run individual services
run-api: build ## Run api server
	@echo "Running api server..."
	@set -a && [ -f ./.env ] && . ./.env && set +a && ./bin/api
	

# docker helpers
docker-run: ## Run docker dev services (attached)
	@docker compose -f docker-compose.dev.yaml up

docker-build: ## Build docker images for services
	@docker compose -f docker-compose.dev.yaml build