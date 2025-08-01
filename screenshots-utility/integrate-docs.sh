#!/bin/bash

# 📝 Documentation Integration Script
# Updates documentation files with generated screenshot links

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}📝 Integrating Screenshots into Documentation${NC}"
echo "=============================================="

# Check if screenshots exist
if [ ! -d "$SCRIPT_DIR/screenshots" ] || [ -z "$(ls -A "$SCRIPT_DIR/screenshots"/*.gif 2>/dev/null)" ]; then
    echo -e "${YELLOW}⚠️ No screenshots found. Run './generate.sh' first.${NC}"
    exit 1
fi

screenshot_count=$(ls -1 "$SCRIPT_DIR/screenshots"/*.gif 2>/dev/null | wc -l)
echo -e "${GREEN}✅ Found $screenshot_count screenshots${NC}"

# =============================================================================
# 📖 Update README.md
# =============================================================================
echo -e "\n${YELLOW}📖 Updating README.md...${NC}"

readme_file="$ROOT_DIR/README.md"
readme_backup="$ROOT_DIR/README.md.backup"

# Create backup
cp "$readme_file" "$readme_backup"

# Create new README section with screenshots
cat > /tmp/readme_screenshots.md << 'EOF'

## 🎬 Visual Demo

### Quick Start Workflow
![Quickstart Demo](screenshots-utility/screenshots/quickstart.gif)

### Performance Optimization in Action
![Performance Demo](screenshots-utility/screenshots/performance.gif)

*Click on the GIFs to see them in full resolution*

EOF

# Insert screenshots section after "Example Output" section if it exists
if grep -q "## 📊 Example Output" "$readme_file"; then
    # Insert after Example Output section
    awk '
        /## 📊 Example Output/ {
            print
            # Print the existing Example Output section until next ## section
            while ((getline line) > 0 && line !~ /^## /) {
                print line
            }
            # Now insert our screenshots section
            system("cat /tmp/readme_screenshots.md")
            # Print the line that starts the next section
            if (line ~ /^## /) print line
            next
        }
        {print}
    ' "$readme_file" > /tmp/readme_new.md
    mv /tmp/readme_new.md "$readme_file"
else
    # Insert before "## 🔧 Installation Options" if it exists
    if grep -q "## 🔧 Installation Options" "$readme_file"; then
        awk '
            /## 🔧 Installation Options/ {
                system("cat /tmp/readme_screenshots.md")
                print
                next
            }
            {print}
        ' "$readme_file" > /tmp/readme_new.md
        mv /tmp/readme_new.md "$readme_file"
    fi
fi

rm -f /tmp/readme_screenshots.md

echo -e "${GREEN}✅ README.md updated${NC}"

# =============================================================================
# ⚡ Update QUICKSTART.md
# =============================================================================
echo -e "\n${YELLOW}⚡ Updating QUICKSTART.md...${NC}"

quickstart_file="$ROOT_DIR/QUICKSTART.md"
quickstart_backup="$ROOT_DIR/QUICKSTART.md.backup"

# Create backup
cp "$quickstart_file" "$quickstart_backup"

# Create quickstart screenshots section
cat > /tmp/quickstart_screenshots.md << 'EOF'

## 🎬 Visual Examples

### Safe Preview (Always Start Here!)
![Dry Run Demo](screenshots-utility/screenshots/dry-run.gif)

### Basic Cleanup
![Basic Usage Demo](screenshots-utility/screenshots/basic-usage.gif)

### Interactive TUI Mode
![TUI Demo](screenshots-utility/screenshots/ui-demo.gif)

*These GIFs show the actual terminal output you'll see*

EOF

# Insert after "What You'll See" section
if grep -q "## 📊 What You'll See" "$quickstart_file"; then
    awk '
        /## 📊 What You.*ll See/ {
            print
            # Print the existing section until next ## section  
            while ((getline line) > 0 && line !~ /^## /) {
                print line
            }
            # Insert screenshots
            system("cat /tmp/quickstart_screenshots.md")
            # Print the next section header
            if (line ~ /^## /) print line
            next
        }
        {print}
    ' "$quickstart_file" > /tmp/quickstart_new.md
    mv /tmp/quickstart_new.md "$quickstart_file"
fi

rm -f /tmp/quickstart_screenshots.md

echo -e "${GREEN}✅ QUICKSTART.md updated${NC}"

# =============================================================================
# 📖 Update USAGE.md
# =============================================================================
echo -e "\n${YELLOW}📖 Updating USAGE.md...${NC}"

usage_file="$ROOT_DIR/USAGE.md"
usage_backup="$ROOT_DIR/USAGE.md.backup"

# Create backup
cp "$usage_file" "$usage_backup"

# Create usage screenshots section
cat > /tmp/usage_screenshots.md << 'EOF'

## 🎬 Visual Examples

### Verbose Mode with Optimization
![Verbose Mode Demo](screenshots-utility/screenshots/verbose.gif)

### Interactive Per-Project Confirmation  
![Interactive Mode Demo](screenshots-utility/screenshots/interactive.gif)

### Performance Optimization Demonstration
![Performance Demo](screenshots-utility/screenshots/performance.gif)

*These demos show the actual CLI experience across different usage modes*

EOF

# Insert after "Usage Modes" section
if grep -q "## 🎛️ Usage Modes" "$usage_file"; then
    awk '
        /## 🎛️ Usage Modes/ {
            print
            # Print the existing section content
            while ((getline line) > 0 && line !~ /^## / && line !~ /^### 1\./) {
                print line
            }
            # Insert screenshots before the mode details
            system("cat /tmp/usage_screenshots.md")
            # Print the line that starts mode details
            if (line ~ /^## / || line ~ /^### /) print line  
            next
        }
        {print}
    ' "$usage_file" > /tmp/usage_new.md
    mv /tmp/usage_new.md "$usage_file"
fi

rm -f /tmp/usage_screenshots.md

echo -e "${GREEN}✅ USAGE.md updated${NC}"

# =============================================================================
# 📊 Summary
# =============================================================================
echo -e "\n${GREEN}🎉 Documentation integration complete!${NC}"
echo "=============================================="

echo -e "📝 Updated files:"
echo "   • README.md (performance + quickstart demos)"
echo "   • QUICKSTART.md (visual examples for key workflows)"  
echo "   • USAGE.md (detailed mode demonstrations)"
echo ""

echo -e "💾 Backups created:"
echo "   • README.md.backup"
echo "   • QUICKSTART.md.backup"
echo "   • USAGE.md.backup"
echo ""

echo -e "${BLUE}💡 Next steps:${NC}"
echo "   1. Review the updated documentation files"
echo "   2. Test the screenshot links in a markdown viewer"
echo "   3. Commit the changes when satisfied"
echo ""

echo -e "${YELLOW}🔧 To revert changes if needed:${NC}"
echo "   mv README.md.backup README.md"
echo "   mv QUICKSTART.md.backup QUICKSTART.md"  
echo "   mv USAGE.md.backup USAGE.md"