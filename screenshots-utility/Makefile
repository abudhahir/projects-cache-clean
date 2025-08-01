# 📸 Screenshot Generation Makefile

# Colors
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m

.PHONY: help install setup screenshots test clean validate benchmark all

help: ## Show this help message
	@echo "📸 Cache Remover Screenshot Utility"
	@echo "=================================="
	@echo ""  
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}'

install: ## Install VHS (required for screenshot generation)
	@echo "$(YELLOW)📦 Installing VHS...$(NC)"
	go install github.com/charmbracelet/vhs@latest
	@echo "$(GREEN)✅ VHS installed$(NC)"

setup: ## Setup test data for screenshots
	@echo "$(YELLOW)🏗️ Setting up test data...$(NC)"
	./setup-test-data.sh
	@echo "$(GREEN)✅ Test data ready$(NC)"

screenshots: ## Generate all screenshots
	@echo "$(YELLOW)📸 Generating all screenshots...$(NC)"
	./generate.sh
	@echo "$(GREEN)✅ Screenshots generated$(NC)"

test: ## Run screenshot generation tests
	@echo "$(YELLOW)🧪 Running screenshot tests...$(NC)"
	cd .. && GENERATE_SCREENSHOTS=1 go test -v ./screenshots-utility/
	@echo "$(GREEN)✅ Tests completed$(NC)"

benchmark: ## Benchmark screenshot generation performance
	@echo "$(YELLOW)⏱️ Running screenshot benchmarks...$(NC)"
	cd .. && BENCHMARK_SCREENSHOTS=1 go test -v -bench=. ./screenshots-utility/
	@echo "$(GREEN)✅ Benchmarks completed$(NC)"

validate: ## Validate generated screenshots
	@echo "$(YELLOW)🔍 Validating screenshots...$(NC)"
	@if [ ! -d "screenshots" ]; then \
		echo "$(RED)❌ Screenshots directory not found. Run 'make screenshots' first.$(NC)"; \
		exit 1; \
	fi
	@count=$$(ls -1 screenshots/*.gif 2>/dev/null | wc -l); \
	if [ $$count -eq 0 ]; then \
		echo "$(RED)❌ No screenshot files found$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Found $$count screenshot files$(NC)"
	@du -sh screenshots/
	@ls -la screenshots/*.gif

clean: ## Clean generated screenshots and test data
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	rm -rf screenshots/  
	rm -rf ../demo-projects/
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

quick: setup screenshots validate ## Quick workflow: setup + generate + validate
	@echo "$(GREEN)🚀 Quick screenshot generation completed!$(NC)"

all: install setup screenshots test validate ## Full workflow: install + setup + generate + test + validate
	@echo "$(GREEN)🎉 Full screenshot generation workflow completed!$(NC)"

# Individual screenshot targets
basic-usage: ## Generate basic usage screenshot only
	@echo "$(YELLOW)📸 Generating basic usage screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/basic-usage.tape

dry-run: ## Generate dry-run screenshot only
	@echo "$(YELLOW)📸 Generating dry-run screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/dry-run.tape

verbose: ## Generate verbose mode screenshot only
	@echo "$(YELLOW)📸 Generating verbose screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/verbose.tape

interactive: ## Generate interactive mode screenshot only
	@echo "$(YELLOW)📸 Generating interactive screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/interactive.tape

ui-demo: ## Generate TUI demo screenshot only
	@echo "$(YELLOW)📸 Generating TUI demo screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/ui-demo.tape

performance: ## Generate performance demo screenshot only
	@echo "$(YELLOW)📸 Generating performance screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/performance.tape

quickstart: ## Generate quickstart workflow screenshot only
	@echo "$(YELLOW)📸 Generating quickstart screenshot...$(NC)"
	cd .. && vhs screenshots-utility/tapes/quickstart.tape

# Development targets
dev-setup: ## Setup for development (installs VHS + builds binary)
	@echo "$(YELLOW)🔧 Setting up development environment...$(NC)"
	@$(MAKE) install
	@cd .. && go build -o cache-remover
	@$(MAKE) setup
	@echo "$(GREEN)✅ Development environment ready$(NC)"

watch: ## Watch for changes and regenerate screenshots (requires entr)
	@echo "$(YELLOW)👀 Watching for changes...$(NC)"
	@if ! command -v entr >/dev/null 2>&1; then \
		echo "$(RED)❌ 'entr' is required for watch mode. Install with: brew install entr$(NC)"; \
		exit 1; \
	fi
	find tapes/ -name "*.tape" | entr -r make screenshots

# CI/CD targets  
ci-test: ## Run tests suitable for CI environment
	@echo "$(YELLOW)🤖 Running CI tests...$(NC)"
	@if ! command -v vhs >/dev/null 2>&1; then \
		echo "$(YELLOW)⚠️ VHS not found, installing...$(NC)"; \
		$(MAKE) install; \
	fi
	@cd .. && go build -o cache-remover
	@$(MAKE) setup
	@cd .. && GENERATE_SCREENSHOTS=1 go test -v -timeout 10m ./screenshots-utility/

# Information targets
info: ## Show system information and requirements
	@echo "$(BLUE)📋 System Information$(NC)"
	@echo "===================="
	@echo "Go version: $$(go version 2>/dev/null || echo 'Not installed')"
	@echo "VHS version: $$(vhs --version 2>/dev/null || echo 'Not installed')"
	@echo "Cache remover: $$([ -f ../cache-remover ] && echo 'Built' || echo 'Not built')"
	@echo "Test data: $$([ -d ../demo-projects ] && echo 'Ready' || echo 'Not setup')"
	@echo "Screenshots: $$([ -d screenshots ] && ls -1 screenshots/*.gif 2>/dev/null | wc -l || echo '0') files"

size: ## Show screenshot file sizes
	@if [ -d "screenshots" ]; then \
		echo "$(BLUE)📊 Screenshot Sizes$(NC)"; \
		echo "=================="; \
		du -h screenshots/*.gif 2>/dev/null | sort -hr || echo "No screenshots found"; \
		echo ""; \
		echo "Total: $$(du -sh screenshots/ 2>/dev/null | cut -f1)"; \
	else \
		echo "$(YELLOW)⚠️ No screenshots directory found$(NC)"; \
	fi