#!/bin/bash
# Install Go and build the cache remover utility

set -e

echo "🔧 Setting up Go Cache Remover Utility"
echo "======================================"

# Check if Go is installed
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "✅ Found Go: $GO_VERSION"
else
    echo "❌ Go is not installed"
    echo ""
    echo "📥 Installing Go..."
    
    # Detect OS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo "Using Homebrew to install Go..."
            brew install go
        else
            echo "Please install Go manually from: https://golang.org/dl/"
            echo "Or install Homebrew first: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        echo "Please install Go using your package manager:"
        echo "  Ubuntu/Debian: sudo apt update && sudo apt install golang-go"
        echo "  CentOS/RHEL:   sudo yum install golang"
        echo "  Arch:          sudo pacman -S go"
        echo "Or download from: https://golang.org/dl/"
        exit 1
    else
        echo "Please install Go from: https://golang.org/dl/"
        exit 1
    fi
fi

# Verify Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go installation failed or not in PATH"
    exit 1
fi

echo "📦 Downloading Go dependencies..."
go mod tidy

echo "🔨 Building cache remover..."
go build -o cache-remover main.go interactive.go webui.go

if [ -f "cache-remover" ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Usage:"
    echo "  ./cache-remover --ui          # Interactive terminal UI"
    echo "  ./cache-remover --web         # Web browser UI"
    echo "  ./cache-remover --help        # Show all options"
    echo ""
    echo "🎯 Quick test:"
    echo "  ./cache-remover --dry-run --verbose"
else
    echo "❌ Build failed"
    exit 1
fi