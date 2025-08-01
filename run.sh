#!/bin/bash

# 🚀 Cache Remover Utility - Quick Run Script
# This script builds (if needed) and runs the cache remover utility

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Cache Remover Utility - Quick Start${NC}"
echo ""

# Check if binary exists and is up to date
REBUILD=false
if [ ! -f "cache-remover-utility" ]; then
    echo -e "${YELLOW}📦 Binary not found, building...${NC}"
    REBUILD=true
elif [ "main.go" -nt "cache-remover-utility" ] || [ "interactive.go" -nt "cache-remover-utility" ] || [ "config.go" -nt "cache-remover-utility" ]; then
    echo -e "${YELLOW}📦 Source code updated, rebuilding...${NC}"
    REBUILD=true
fi

# Build if needed
if [ "$REBUILD" = true ]; then
    ./build.sh
    echo ""
fi

# Parse command line arguments
UI_MODE=false
DRY_RUN=false
DIRECTORY="."

# Simple argument parsing
for arg in "$@"; do
    case $arg in
        --ui|-ui)
            UI_MODE=true
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help|-h)
            ./cache-remover-utility --help
            exit 0
            ;;
        -*)
            # Pass through other flags
            ;;
        *)
            # Assume it's a directory
            if [ -d "$arg" ]; then
                DIRECTORY="$arg"
            fi
            ;;
    esac
done

# Show what we're about to do
echo -e "${GREEN}🎯 Running Cache Remover:${NC}"
if [ "$UI_MODE" = true ]; then
    echo -e "   Mode: ${BLUE}Interactive TUI${NC}"
elif [ "$DRY_RUN" = true ]; then
    echo -e "   Mode: ${YELLOW}Dry Run (Preview Only)${NC}"
else
    echo -e "   Mode: ${GREEN}Command Line${NC}"
fi
echo -e "   Directory: ${BLUE}$DIRECTORY${NC}"
echo ""

# Run the application
if [ $# -eq 0 ]; then
    # No arguments provided, show menu
    echo -e "${YELLOW}No arguments provided. Choose an option:${NC}"
    echo ""
    echo "1) 🖥️  Interactive TUI (recommended)"
    echo "2) 🔍 Dry run preview"
    echo "3) 📋 Show help"
    echo "4) 💻 Command line cleanup (with confirmation)"
    echo ""
    read -p "Enter choice (1-4): " choice
    
    case $choice in
        1)
            exec ./cache-remover-utility --ui "$DIRECTORY"
            ;;
        2)
            exec ./cache-remover-utility --dry-run "$DIRECTORY"
            ;;
        3)
            exec ./cache-remover-utility --help
            ;;
        4)
            exec ./cache-remover-utility "$DIRECTORY"
            ;;
        *)
            echo -e "${RED}Invalid choice. Running TUI by default...${NC}"
            exec ./cache-remover-utility --ui "$DIRECTORY"
            ;;
    esac
else
    # Arguments provided, run directly
    exec ./cache-remover-utility "$@"
fi