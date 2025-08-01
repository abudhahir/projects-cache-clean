#!/bin/bash

# 🔨 Cache Remover Utility - Build Script
# This script builds the cache remover utility with proper optimizations

set -e  # Exit on any error

echo "🔨 Building Cache Remover Utility..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first:"
    echo "   Run: ./install-go.sh"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
echo "✅ Using Go $GO_VERSION"

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -f cache-remover cache-remover-utility

# Build with optimizations
echo "🚀 Building optimized binary..."
go build -ldflags="-s -w" -o cache-remover-utility .

# Verify the build
if [ -f "cache-remover-utility" ]; then
    echo "✅ Build successful!"
    
    # Show binary info
    echo ""
    echo "📊 Binary Information:"
    ls -lh cache-remover-utility
    
    # Test basic functionality
    echo ""
    echo "🧪 Testing basic functionality:"
    ./cache-remover-utility --help | head -3
    
    echo ""
    echo "🎉 Build complete! You can now run:"
    echo "   ./cache-remover-utility --help"
    echo "   ./cache-remover-utility --ui"
    echo "   ./cache-remover-utility --dry-run ."
else
    echo "❌ Build failed!"
    exit 1
fi