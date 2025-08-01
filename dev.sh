#!/bin/bash

# 🛠️ Cache Remover Utility - Development Script
# This script provides common development tasks

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

show_help() {
    echo -e "${BLUE}🛠️  Cache Remover Utility - Development Script${NC}"
    echo ""
    echo "Usage: ./dev.sh [command]"
    echo ""
    echo "Commands:"
    echo -e "  ${GREEN}build${NC}      - Build the application"
    echo -e "  ${GREEN}test${NC}       - Run all tests"
    echo -e "  ${GREEN}test-race${NC}  - Run tests with race detection"
    echo -e "  ${GREEN}lint${NC}       - Run code linting (gofmt, go vet)"
    echo -e "  ${GREEN}clean${NC}      - Clean build artifacts"
    echo -e "  ${GREEN}run${NC}        - Build and run with TUI"
    echo -e "  ${GREEN}demo${NC}       - Run with demo data"
    echo -e "  ${GREEN}install${NC}    - Install to /usr/local/bin"
    echo -e "  ${GREEN}deps${NC}       - Download and verify dependencies"
    echo -e "  ${GREEN}help${NC}       - Show this help"
    echo ""
    echo "Examples:"
    echo "  ./dev.sh build     # Build the application"
    echo "  ./dev.sh test      # Run tests"
    echo "  ./dev.sh run --ui  # Build and run TUI"
}

build_app() {
    echo -e "${BLUE}🔨 Building application...${NC}"
    ./build.sh
}

run_tests() {
    echo -e "${BLUE}🧪 Running tests...${NC}"
    go test -v
}

run_tests_race() {
    echo -e "${BLUE}🏃 Running tests with race detection...${NC}"
    go test -race -v
}

run_lint() {
    echo -e "${BLUE}🧹 Running linter...${NC}"
    
    echo "Checking formatting..."
    if ! gofmt -l . | grep -q .; then
        echo -e "${GREEN}✅ Code is properly formatted${NC}"
    else
        echo -e "${RED}❌ Code needs formatting:${NC}"
        gofmt -l .
        echo "Run: gofmt -w ."
        exit 1
    fi
    
    echo "Running go vet..."
    if go vet ./...; then
        echo -e "${GREEN}✅ No issues found by go vet${NC}"
    else
        echo -e "${RED}❌ Issues found by go vet${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ All linting checks passed${NC}"
}

clean_build() {
    echo -e "${BLUE}🧹 Cleaning build artifacts...${NC}"
    rm -f cache-remover cache-remover-utility
    rm -rf test-area tree-test
    echo -e "${GREEN}✅ Clean complete${NC}"
}

run_app() {
    echo -e "${BLUE}🚀 Building and running application...${NC}"
    build_app
    echo ""
    shift # Remove 'run' from arguments
    ./cache-remover-utility "$@"
}

run_demo() {
    echo -e "${BLUE}🎬 Setting up demo environment...${NC}"
    
    # Create demo projects if they don't exist
    if [ ! -d "demo-projects" ]; then
        mkdir -p demo-projects/{frontend/{react-app,vue-app},backend/{java-service,python-api},mobile/flutter-app}
        
        # Frontend
        echo '{"name":"react-app","version":"1.0.0"}' > demo-projects/frontend/react-app/package.json
        mkdir -p demo-projects/frontend/react-app/{node_modules,dist,build}
        echo "demo content" > demo-projects/frontend/react-app/node_modules/react.js
        
        echo '{"name":"vue-app","version":"2.0.0"}' > demo-projects/frontend/vue-app/package.json
        mkdir -p demo-projects/frontend/vue-app/{node_modules,dist}
        echo "vue content" > demo-projects/frontend/vue-app/node_modules/vue.js
        
        # Backend
        echo '<project></project>' > demo-projects/backend/java-service/pom.xml
        mkdir -p demo-projects/backend/java-service/target/classes
        echo "compiled class" > demo-projects/backend/java-service/target/classes/App.class
        
        echo "flask==2.0.1" > demo-projects/backend/python-api/requirements.txt
        mkdir -p demo-projects/backend/python-api/{__pycache__,build}
        echo "compiled python" > demo-projects/backend/python-api/__pycache__/app.pyc
        
        # Mobile
        echo 'name: flutter_app' > demo-projects/mobile/flutter-app/pubspec.yaml
        mkdir -p demo-projects/mobile/flutter-app/{build,.dart_tool}
        echo "flutter build" > demo-projects/mobile/flutter-app/build/app.so
        
        echo -e "${GREEN}✅ Demo projects created${NC}"
    fi
    
    build_app
    echo ""
    echo -e "${YELLOW}🎬 Running with demo projects (tree view)...${NC}"
    ./cache-remover-utility --ui demo-projects/
}

install_app() {
    echo -e "${BLUE}📦 Installing application...${NC}"
    build_app
    
    if [ ! -f "cache-remover-utility" ]; then
        echo -e "${RED}❌ Build failed, cannot install${NC}"
        exit 1
    fi
    
    echo "Installing to /usr/local/bin..."
    sudo cp cache-remover-utility /usr/local/bin/
    sudo chmod +x /usr/local/bin/cache-remover-utility
    
    echo -e "${GREEN}✅ Installed successfully!${NC}"
    echo "You can now run: cache-remover-utility --ui"
}

download_deps() {
    echo -e "${BLUE}📥 Downloading dependencies...${NC}"
    go mod download
    go mod verify
    echo -e "${GREEN}✅ Dependencies verified${NC}"
}

# Main script logic
case "${1:-help}" in
    build)
        build_app
        ;;
    test)
        run_tests
        ;;
    test-race)
        run_tests_race
        ;;
    lint)
        run_lint
        ;;
    clean)
        clean_build
        ;;
    run)
        run_app "$@"
        ;;
    demo)
        run_demo
        ;;
    install)
        install_app
        ;;
    deps)
        download_deps
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac