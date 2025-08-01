#!/bin/bash
# Installation script for Cache Remover Utility

set -e

echo "🧹 Cache Remover Utility Installation"
echo "======================================"

# Check Python version
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is required but not installed."
    echo "Please install Python 3.8+ and try again."
    exit 1
fi

PYTHON_VERSION=$(python3 -c 'import sys; print(".".join(map(str, sys.version_info[:2])))')
echo "✅ Found Python $PYTHON_VERSION"

# Make script executable
chmod +x cache_remover.py
echo "✅ Made cache_remover.py executable"

# Create symlink in /usr/local/bin if it exists and is writable
if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    ln -sf "$(pwd)/cache_remover.py" /usr/local/bin/cache-remover
    echo "✅ Created symlink: /usr/local/bin/cache-remover"
    echo ""
    echo "🎉 Installation complete!"
    echo "You can now run 'cache-remover' from anywhere."
else
    echo "⚠️  Could not create system-wide symlink."
    echo "You can run the utility with: python3 $(pwd)/cache_remover.py"
fi

echo ""
echo "Usage examples:"
echo "  cache-remover --help"
echo "  cache-remover --dry-run --verbose"
echo "  cache-remover --dir ~/Projects --workers 16"