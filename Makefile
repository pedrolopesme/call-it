# 🚀 Call-It - Make HTTP calls like a boss!
# ================================================

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary info
BINARY_NAME=call-it
BUILD_DIR=./bin
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)

# Build info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell $(GOCMD) version | cut -d' ' -f3)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[0;37m
BOLD=\033[1m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## Show this help message
help:
	@echo ""
	@echo "$(BOLD)$(CYAN)🚀 Call-It Makefile$(NC)"
	@echo "$(PURPLE)================================================$(NC)"
	@echo ""
	@echo "$(BOLD)Available targets:$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  $(CYAN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)Project Info:$(NC)"
	@echo "  $(YELLOW)Version:$(NC)     $(VERSION)"
	@echo "  $(YELLOW)Go Version:$(NC)  $(GO_VERSION)"
	@echo "  $(YELLOW)Binary:$(NC)      $(BINARY_PATH)"
	@echo ""

build: ## Build the application
	@echo "$(BOLD)$(GREEN)🔨 Building call-it...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "  $(BLUE)→$(NC) Version: $(VERSION)"
	@echo "  $(BLUE)→$(NC) Build time: $(BUILD_TIME)"
	@$(GOBUILD) -ldflags "-X github.com/pedrolopesme/call-it/internal/version.Version=$(VERSION) -X github.com/pedrolopesme/call-it/internal/version.BuildTime=$(BUILD_TIME)" -o $(BINARY_PATH) -v ./cmd/call-it
	@echo "$(GREEN)✅ Build complete!$(NC) Binary: $(BINARY_PATH)"

test: ## Run all tests
	@echo "$(BOLD)$(YELLOW)🧪 Running tests...$(NC)"
	@$(GOTEST) -v -race -cover ./...
	@echo "$(GREEN)✅ Tests completed!$(NC)"

test-coverage: ## Run tests with coverage report
	@echo "$(BOLD)$(YELLOW)📊 Running tests with coverage...$(NC)"
	@$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated!$(NC) Open coverage.html in your browser"

clean: ## Clean build artifacts
	@echo "$(BOLD)$(RED)🧹 Cleaning up...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Clean complete!$(NC)"

run: build ## Build and run the application
	@echo "$(BOLD)$(GREEN)🚀 Running call-it...$(NC)"
	@echo "  $(BLUE)→$(NC) Args: $(filter-out $@,$(MAKECMDGOALS))"
	@$(BINARY_PATH) $(filter-out $@,$(MAKECMDGOALS))

fmt: ## Format Go code
	@echo "$(BOLD)$(CYAN)📝 Formatting code...$(NC)"
	@$(GOFMT) -w .
	@echo "$(GREEN)✅ Code formatted!$(NC)"

lint: ## Lint the code
	@echo "$(BOLD)$(PURPLE)🔍 Linting code...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Linting complete!$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not found. Install it with:$(NC)"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Vet the code
	@echo "$(BOLD)$(BLUE)🔍 Vetting code...$(NC)"
	@$(GOVET) ./...
	@echo "$(GREEN)✅ Vet complete!$(NC)"

deps: ## Update dependencies
	@echo "$(BOLD)$(CYAN)📦 Updating dependencies...$(NC)"
	@$(GOMOD) tidy
	@$(GOMOD) download
	@echo "$(GREEN)✅ Dependencies updated!$(NC)"

deps-check: ## Check for outdated dependencies
	@echo "$(BOLD)$(YELLOW)📋 Checking for outdated dependencies...$(NC)"
	@if command -v go-mod-outdated >/dev/null 2>&1; then \
		$(GOCMD) list -u -m -json all | go-mod-outdated -update -direct; \
	else \
		echo "$(YELLOW)⚠️  go-mod-outdated not found. Install it with:$(NC)"; \
		echo "  go install github.com/psampaz/go-mod-outdated@latest"; \
	fi

install: build ## Install the binary to GOPATH/bin
	@echo "$(BOLD)$(GREEN)📦 Installing call-it...$(NC)"
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

info: ## Show project information
	@echo ""
	@echo "$(BOLD)$(CYAN)📋 Project Information$(NC)"
	@echo "$(PURPLE)================================================$(NC)"
	@echo "$(YELLOW)Project:$(NC)     call-it"
	@echo "$(YELLOW)Version:$(NC)     $(VERSION)"
	@echo "$(YELLOW)Go Version:$(NC)  $(GO_VERSION)"
	@echo "$(YELLOW)Build Time:$(NC)  $(BUILD_TIME)"
	@echo "$(YELLOW)Binary:$(NC)      $(BINARY_PATH)"
	@echo "$(YELLOW)Build Dir:$(NC)   $(BUILD_DIR)"
	@echo ""

ci: fmt vet test build ## Run full CI pipeline (fmt, vet, test, build)
	@echo "$(BOLD)$(GREEN)🎉 CI pipeline completed successfully!$(NC)"

.PHONY: help build test test-coverage clean run fmt lint vet deps deps-check install info ci