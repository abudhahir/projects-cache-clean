# 🔨 Cache Remover Utility - Makefile
# Simple Makefile for common development tasks

.PHONY: build test clean run install help demo lint deps

# Default target
all: build

# Build the application
build:
	@echo "🔨 Building Cache Remover Utility..."
	@./build.sh

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v

# Run tests with race detection
test-race:
	@echo "🏃 Running tests with race detection..."
	@go test -race -v

# Run linting checks
lint:
	@echo "🧹 Running linter..."
	@./dev.sh lint

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f cache-remover cache-remover-utility
	@rm -rf test-area tree-test demo-projects

# Build and run with TUI
run: build
	@echo "🚀 Running Cache Remover TUI..."
	@./cache-remover-utility --ui

# Run with demo data
demo:
	@echo "🎬 Running demo..."
	@./dev.sh demo

# Install to system
install: build
	@echo "📦 Installing to /usr/local/bin..."
	@sudo cp cache-remover-utility /usr/local/bin/
	@sudo chmod +x /usr/local/bin/cache-remover-utility
	@echo "✅ Installed successfully!"

# Download and verify dependencies
deps:
	@echo "📥 Downloading dependencies..."
	@go mod download
	@go mod verify

# Quick start with current directory
start: build
	@./cache-remover-utility --ui .

# Dry run with current directory
preview: build
	@./cache-remover-utility --dry-run .

# Show help
help:
	@echo "🛠️  Cache Remover Utility - Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  test      - Run all tests"
	@echo "  test-race - Run tests with race detection"
	@echo "  lint      - Run code linting"
	@echo "  clean     - Clean build artifacts"
	@echo "  run       - Build and run with TUI"
	@echo "  demo      - Run with demo data"
	@echo "  start     - Quick start with current directory"
	@echo "  preview   - Dry run with current directory"
	@echo "  install   - Install to /usr/local/bin"
	@echo "  deps      - Download dependencies"
	@echo "  help      - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build    # Build the application"
	@echo "  make run      # Build and run TUI"
	@echo "  make demo     # Run with demo projects"