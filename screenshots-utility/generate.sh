#!/bin/bash

# 📸 Screenshot Generation Script
# Generates all documentation screenshots using VHS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}📸 Cache Remover Screenshot Generation${NC}"
echo "=================================================="

# =============================================================================
# 🔍 Pre-flight Checks
# =============================================================================
echo -e "\n${YELLOW}🔍 Running pre-flight checks...${NC}"

# Check if VHS is installed, install if missing
if ! command -v vhs &> /dev/null; then
    echo -e "${YELLOW}📦 VHS not found, installing automatically...${NC}"
    if command -v go &> /dev/null; then
        go install github.com/charmbracelet/vhs@latest
        # Add ~/go/bin to PATH if it's not already there
        export PATH="$HOME/go/bin:$PATH"
        if command -v vhs &> /dev/null; then
            echo -e "${GREEN}✅ VHS installed successfully${NC}"
        else
            echo -e "${RED}❌ VHS installation failed. Please install manually:${NC}"
            echo "   go install github.com/charmbracelet/vhs@latest"
            echo "   export PATH=\"\$HOME/go/bin:\$PATH\""
            exit 1
        fi
    else
        echo -e "${RED}❌ Go is not installed. Please install Go first:${NC}"
        echo "   brew install go  # macOS"
        echo "   apt install golang-go  # Ubuntu"
        exit 1
    fi
else
    echo -e "${GREEN}✅ VHS found: $(vhs --version)${NC}"
fi

# Check if ttyd is installed, install if missing (required by VHS)
if ! command -v ttyd &> /dev/null; then
    echo -e "${YELLOW}📦 ttyd not found, installing automatically...${NC}"
    if command -v brew &> /dev/null; then
        brew install ttyd
        if command -v ttyd &> /dev/null; then
            echo -e "${GREEN}✅ ttyd installed successfully${NC}"
        else
            echo -e "${RED}❌ ttyd installation failed${NC}"
            exit 1
        fi
    elif command -v apt &> /dev/null; then
        sudo apt update && sudo apt install -y ttyd
        echo -e "${GREEN}✅ ttyd installed via apt${NC}"
    else
        echo -e "${RED}❌ Cannot auto-install ttyd. Please install manually:${NC}"
        echo "   brew install ttyd  # macOS"
        echo "   apt install ttyd   # Ubuntu"
        echo "   Or visit: https://github.com/tsl0922/ttyd"
        exit 1
    fi
else
    echo -e "${GREEN}✅ ttyd found${NC}"
fi

# Check if cache-remover is built
if [ ! -f "$ROOT_DIR/cache-remover" ]; then
    echo -e "${YELLOW}📦 Building cache-remover...${NC}"
    cd "$ROOT_DIR"
    go build -o cache-remover
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ cache-remover built successfully${NC}"
    else
        echo -e "${RED}❌ Failed to build cache-remover${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ cache-remover binary found${NC}"
fi

# =============================================================================
# 🏗️ Setup Test Data
# =============================================================================
echo -e "\n${YELLOW}🏗️ Setting up test data...${NC}"
cd "$SCRIPT_DIR"
./setup-test-data.sh

# Move to root directory for recording
cd "$ROOT_DIR"

# =============================================================================
# 📁 Create Screenshots Directory
# =============================================================================
echo -e "\n${YELLOW}📁 Preparing screenshots directory...${NC}"
mkdir -p "$SCRIPT_DIR/screenshots"

# =============================================================================
# 🎬 Generate Screenshots
# =============================================================================
screenshots=(
    "basic-usage:Basic Usage Demo"
    "dry-run:Dry Run Demo" 
    "verbose:Verbose Mode Demo"
    "interactive:Interactive Mode Demo"
    "ui-demo:TUI Interface Demo"
    "performance:Performance Optimization Demo"
    "quickstart:Quickstart Workflow Demo"
)

echo -e "\n${YELLOW}🎬 Generating screenshots...${NC}"
total=${#screenshots[@]}
current=0

for screenshot in "${screenshots[@]}"; do
    IFS=':' read -r name description <<< "$screenshot"
    current=$((current + 1))
    
    echo -e "\n${BLUE}[$current/$total] Generating $description...${NC}"
    
    tape_file="$SCRIPT_DIR/tapes/${name}.tape"
    if [ ! -f "$tape_file" ]; then
        echo -e "${RED}❌ Tape file not found: $tape_file${NC}"
        continue
    fi
    
    # Run VHS to generate screenshot
    if vhs "$tape_file"; then
        output_file="$SCRIPT_DIR/screenshots/${name}.gif"
        if [ -f "$output_file" ]; then
            size=$(du -h "$output_file" | cut -f1)
            echo -e "${GREEN}✅ Generated $name.gif ($size)${NC}"
        else
            echo -e "${RED}❌ Expected output file not found: $output_file${NC}"
        fi
    else
        echo -e "${RED}❌ Failed to generate $name.gif${NC}"
    fi
done

# =============================================================================
# 📊 Generate Summary Report
# =============================================================================
echo -e "\n${YELLOW}📊 Generating summary report...${NC}"

report_file="$SCRIPT_DIR/screenshots/README.md"
cat > "$report_file" << 'EOF'
# 📸 Generated Screenshots

This directory contains all generated screenshots for the Cache Remover documentation.

## 🎯 Screenshot Index

| Screenshot | Purpose | Size | Documentation |
|------------|---------|------|---------------|
EOF

# Add screenshots to report
for screenshot in "${screenshots[@]}"; do
    IFS=':' read -r name description <<< "$screenshot"
    gif_file="$SCRIPT_DIR/screenshots/${name}.gif"
    
    if [ -f "$gif_file" ]; then
        size=$(du -h "$gif_file" | cut -f1)
        
        # Determine which docs use this screenshot
        docs=""
        case $name in
            "basic-usage"|"quickstart")
                docs="README.md, QUICKSTART.md"
                ;;
            "dry-run")
                docs="QUICKSTART.md, USAGE.md"
                ;;
            "verbose"|"interactive")
                docs="USAGE.md"
                ;;
            "ui-demo")
                docs="README.md, QUICKSTART.md"
                ;;
            "performance")
                docs="README.md, USAGE.md"
                ;;
        esac
        
        echo "| \`${name}.gif\` | $description | $size | $docs |" >> "$report_file"
    fi
done

cat >> "$report_file" << 'EOF'

## 🖼️ Preview

### Basic Usage
![Basic Usage](basic-usage.gif)

### Dry Run Mode  
![Dry Run](dry-run.gif)

### Verbose Mode
![Verbose Mode](verbose.gif)

### Interactive Mode
![Interactive Mode](interactive.gif)

### TUI Interface
![TUI Demo](ui-demo.gif)

### Performance Demo
![Performance](performance.gif)

### Quickstart Workflow
![Quickstart](quickstart.gif)

## 📝 Usage in Documentation

### For QUICKSTART.md
```markdown
![Basic Usage](screenshots-utility/screenshots/basic-usage.gif)
![Dry Run](screenshots-utility/screenshots/dry-run.gif)
![TUI Demo](screenshots-utility/screenshots/ui-demo.gif)
```

### For USAGE.md
```markdown
![Verbose Mode](screenshots-utility/screenshots/verbose.gif)
![Interactive Mode](screenshots-utility/screenshots/interactive.gif)
![Performance Demo](screenshots-utility/screenshots/performance.gif)
```

### For README.md
```markdown
![Quickstart](screenshots-utility/screenshots/quickstart.gif)
![Performance Demo](screenshots-utility/screenshots/performance.gif)
```

---
*Generated by screenshots-utility/generate.sh*
EOF

echo -e "${GREEN}✅ Summary report created: $report_file${NC}"

# =============================================================================
# 📊 Final Summary
# =============================================================================
echo -e "\n${GREEN}🎉 Screenshot generation complete!${NC}"
echo "=================================================="

screenshot_count=$(ls -1 "$SCRIPT_DIR/screenshots"/*.gif 2>/dev/null | wc -l || echo 0)
total_size=$(du -sh "$SCRIPT_DIR/screenshots" 2>/dev/null | cut -f1 || echo "0")

echo -e "📁 Output directory: ${BLUE}$SCRIPT_DIR/screenshots${NC}"
echo -e "📸 Screenshots generated: ${GREEN}$screenshot_count${NC}"
echo -e "💾 Total size: ${YELLOW}$total_size${NC}"

if [ "$screenshot_count" -gt 0 ]; then
    echo ""
    echo -e "${YELLOW}📋 Generated files:${NC}"
    ls -la "$SCRIPT_DIR/screenshots"/*.gif 2>/dev/null | while read -r line; do
        echo "   $line"
    done
fi

echo ""
echo -e "${BLUE}💡 Next steps:${NC}"
echo "   1. Review generated screenshots in: screenshots-utility/screenshots/"
echo "   2. Update documentation with new screenshot links"
echo "   3. Commit changes to repository"
echo ""
echo -e "${GREEN}🚀 Screenshots are ready for documentation!${NC}"