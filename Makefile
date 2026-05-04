.PHONY: test test-unit test-integration test-e2e test-all cover build clean

.DEFAULT_GOAL := help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Build

build: ## Build the project
	go build ./...

clean: ## Clean build artifacts
	go clean -testcache ./...

##@ Testing

test-unit: ## Run unit tests
	@echo "🧪 Running unit tests..."
	go test -v ./internal/engine ./internal/tools

test-integration: ## Run integration tests
	@echo "🔧 Running integration tests..."
	go test -v ./integration

test-e2e: ## Run end-to-end tests
	@echo "🚀 Running e2e tests..."
	@bash scripts/test-e2e.sh

test: ## Run all tests
	@echo "🧪 Running all tests..."
	go test -v ./internal/... ./integration

test-all: ## Run all tests with verbose output
	@echo "🧪 Running all tests..."
	go test -v -failfast ./internal/... ./integration

cover: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./internal/engine ./integration
	go tool cover -html=coverage.out
	@echo "📊 Coverage report generated: coverage.html"

##@ Validation

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping"; \
	fi

fmt: ## Format code
	@echo "🎨 Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...

validate: fmt vet lint ## Run all validations (fmt, vet, lint)
