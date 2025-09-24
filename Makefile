.PHONY: help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build web and api binary
	@echo "Building binaries..."
	go build -o bin/api ./backend/


# run individual services
run-api: build ## Run api server
	@echo "Running api server..."
	@set -a && [ -f ./.env ] && . ./.env && set +a && ./bin/api
	

# docker helpers
docker-run: ## Run docker dev services (attached)
	@docker compose -f docker-compose.dev.yaml up

docker-build: ## Build docker images for services
	@docker compose -f docker-compose.dev.yaml build