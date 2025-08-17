.PHONY: help build run dev test clean docker-build docker-run migrate

# Default target
help: ## Show this help message
	@echo "Minimart Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build commands
build: ## Build the application
	$(go env GOPATH)/bin/templ generate
	go build -o bin/minimart ./cmd/server

run: ## Run the application
	$(go env GOPATH)/bin/templ generate
	go run ./cmd/server/main.go

generate: ## Generate Templ templates
	$(go env GOPATH)/bin/templ generate

dev: ## Start development server with hot reload
	@echo "Starting development server..."
	@echo "Templates will auto-reload on changes"
	@which air > /dev/null || (echo "Installing Air..." && go install github.com/air-verse/air@latest)
	air -c .air.toml

# Testing
test: ## Run all tests
	go test ./...

test-verbose: ## Run all tests with verbose output
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -cover ./...

test-integration: ## Run integration tests only
	go test ./internal/integration/...

# Database
migrate: ## Run database migrations
	@echo "Running database migrations..."
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down: ## Rollback last migration
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status: ## Show migration status
	goose -dir migrations postgres "$(DATABASE_URL)" status

# Docker commands
docker-build: ## Build Docker image
	docker build -t minimart:latest .

docker-run: ## Run application in Docker
	docker-compose up --build

docker-dev: ## Run development environment with Docker
	docker-compose up -d
	@echo "Services started. Access app at http://localhost:3000"

docker-clean: ## Clean up Docker containers and images
	docker-compose down
	docker system prune -f

# Utilities
clean: ## Clean build artifacts
	rm -rf bin/ tmp/ 

fmt: ## Format Go code
	go fmt ./...

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

# Template debugging
debug-templates: ## Run with template debugging enabled
	DEBUG=templates go run ./cmd/server/main.go

# Environment setup
setup: ## Set up development environment
	@echo "Setting up development environment..."
	go mod download
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@which goose > /dev/null || go install github.com/pressly/goose/v3/cmd/goose@latest
	cp .env.example .env
	@echo "Development environment ready!"
	@echo "Edit .env file with your configuration"

# Production build
build-prod: ## Build optimized production binary
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/minimart ./cmd/server

# Status check
status: ## Check application and dependency status
	@echo "Go version:"
	@go version
	@echo "\nDependencies:"
	@go list -m all | head -20
	@echo "\nDatabase connection test:"
	@if [ -n "$(DATABASE_URL)" ]; then \
		go run ./cmd/server/main.go --check-db 2>/dev/null && echo "✓ Database OK" || echo "✗ Database connection failed"; \
	else \
		echo "DATABASE_URL not set"; \
	fi
