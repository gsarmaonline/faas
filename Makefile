# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=faas
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
.DEFAULT_GOAL := help

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help
help: ## Display this help message
	@echo "$(BLUE)FAAS Project Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## Run all tests
	@echo "$(YELLOW)Running all tests...$(NC)"
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

.PHONY: test-short
test-short: ## Run tests in short mode
	@echo "$(YELLOW)Running tests in short mode...$(NC)"
	$(GOTEST) -short -v ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "$(YELLOW)Running tests with race detection...$(NC)"
	$(GOTEST) -race -v ./...

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo "$(YELLOW)Running benchmark tests...$(NC)"
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: build
build: ## Build the application
	@echo "$(YELLOW)Building application...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	@echo "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "$(YELLOW)Building for Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./...
	@echo "$(GREEN)Linux build completed: $(BINARY_UNIX)$(NC)"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html
	@echo "$(GREEN)Clean completed$(NC)"

.PHONY: deps
deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GOGET) -d ./...
	@echo "$(GREEN)Dependencies downloaded$(NC)"

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "$(YELLOW)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

.PHONY: mod-tidy
mod-tidy: ## Clean up go.mod and go.sum
	@echo "$(YELLOW)Tidying go.mod...$(NC)"
	$(GOMOD) tidy
	@echo "$(GREEN)go.mod tidied$(NC)"

.PHONY: mod-verify
mod-verify: ## Verify dependencies
	@echo "$(YELLOW)Verifying dependencies...$(NC)"
	$(GOMOD) verify
	@echo "$(GREEN)Dependencies verified$(NC)"

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "$(YELLOW)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)Linting completed$(NC)"; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(NC)"
	$(GOCMD) vet ./...
	@echo "$(GREEN)go vet completed$(NC)"

.PHONY: check
check: fmt vet test ## Run format, vet, and test
	@echo "$(GREEN)All checks passed!$(NC)"

.PHONY: ci
ci: mod-verify check test-race test-coverage ## Run full CI pipeline
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

.PHONY: dev-setup
dev-setup: deps ## Set up development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing golangci-lint...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "$(GREEN)Development environment ready!$(NC)"

.PHONY: run-examples
run-examples: build ## Build and show usage examples
	@echo "$(YELLOW)FAAS Framework built successfully!$(NC)"
	@echo ""
	@echo "$(BLUE)Example usage patterns:$(NC)"
	@echo "$(GREEN)• Run all tests:$(NC) make test"
	@echo "$(GREEN)• Run with coverage:$(NC) make test-coverage"
	@echo "$(GREEN)• Format and check:$(NC) make check"
	@echo "$(GREEN)• Full CI pipeline:$(NC) make ci"

.PHONY: docker-build
docker-build: ## Build Docker image (if Dockerfile exists)
	@if [ -f "Dockerfile" ]; then \
		echo "$(YELLOW)Building Docker image...$(NC)"; \
		docker build -t $(BINARY_NAME):latest .; \
		echo "$(GREEN)Docker image built: $(BINARY_NAME):latest$(NC)"; \
	else \
		echo "$(RED)Dockerfile not found$(NC)"; \
	fi

.PHONY: watch-test
watch-test: ## Watch for changes and run tests (requires entr)
	@echo "$(YELLOW)Watching for changes and running tests...$(NC)"
	@if command -v entr >/dev/null 2>&1; then \
		find . -name "*.go" | entr -c make test; \
	else \
		echo "$(RED)entr not installed. Install with: brew install entr (macOS) or apt-get install entr (Ubuntu)$(NC)"; \
	fi
