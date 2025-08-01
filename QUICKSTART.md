# ⚡ Quick Start Guide

Get up and running with Cache Remover in under 2 minutes!

## 🚀 Installation

```bash
# Clone the repository
git clone https://github.com/abudhahir/projects-cache-clean.git
cd projects-cache-clean

# Option 1: Use build script (recommended)
./build.sh

# Option 2: Use Makefile
make build

# Option 3: Manual build
go build -o cache-remover-utility

# Install Go first if needed
./install-go.sh
```

## 📦 System Installation (Optional)

```bash
# Install to system PATH for global access
make install
# OR
sudo cp cache-remover-utility /usr/local/bin/

# Then use from anywhere:
cache-remover-utility --ui ~/Projects
```

## 🏃 Quick Start Options

```bash
# Option 1: Smart run script with menu
./run.sh

# Option 2: Direct commands
./cache-remover-utility --dry-run ~/Projects    # Safe preview
./cache-remover-utility ~/Projects              # Clean with confirmation  
./cache-remover-utility --ui ~/Projects         # Interactive tree view

# Option 3: Using make
make run        # Build and run TUI
make demo       # Run with demo projects
make preview    # Dry run current directory
```

## 💨 Most Common Usage

```bash
# See what would be cleaned (safe preview)
./cache-remover-utility --dry-run ~/Projects

# Clean everything (with confirmation)
./cache-remover-utility ~/Projects

# Interactive tree view (NEW!)
./cache-remover-utility --ui ~/Projects
```

## 📊 What You'll See

**Dry run example:**
```
🧹 Cache Remover Utility
Scanning directory: /Users/dev/Projects
Workers: 8
🔍 DRY RUN MODE - No files will be removed

📁 Found project: /Users/dev/Projects/my-react-app
⏭️  Skipping cache directory: /Users/dev/Projects/my-react-app/node_modules
Found 3 projects

🗂️  my-react-app (Node.js): 2 cache items (245.2 MB)
🔍 Would remove 2 items (245.2 MB) from: /Users/dev/Projects/my-react-app
  - /Users/dev/Projects/my-react-app/node_modules (242.1 MB)
  - /Users/dev/Projects/my-react-app/build (3.1 MB)

📊 Cleanup Statistics:
   Projects processed: 3
   Cache items removed: 6
   Total space reclaimed: 414.4 MB
   Processing time: 1.2s
```

## 🛠️ Development & Build Commands

| Command | Purpose |
|---------|---------|
| `./build.sh` | **Build** - compile optimized binary |
| `./run.sh` | **Smart run** - build and run with menu |
| `./dev.sh demo` | **Demo** - run with generated test projects |
| `./dev.sh test` | **Test** - run all tests |
| `make install` | **Install** - install to system PATH |

## 🎯 Essential Usage Commands

| Command | Purpose |
|---------|---------|
| `./cache-remover-utility ~/Projects` | **Most common** - scan and clean with confirmation |
| `./cache-remover-utility --dry-run ~/Projects` | **Safe preview** - see what would be removed |
| `./cache-remover-utility --ui ~/Projects` | **Interactive tree** - select with TUI |
| `./cache-remover-utility --verbose ~/Projects` | **Detailed** - see all scanning activity |

## 🛡️ Safety First

- **Always start with `-dry-run`** to see what would be removed
- The tool **only removes rebuildable cache** (node_modules, __pycache__, target, etc.)
- **Asks for confirmation** before actual cleanup
- **Skip projects easily** with interactive mode

## ⚡ Performance Features

- **4,700x faster** cache directory processing with advanced optimization
- **Smart cache boundary detection** - stops at cache directory boundaries instead of scanning contents
- **Depth-limited scanning** - optimized 3-level depth calculation for cache directories
- **Extension search optimization** - skips cache directories during file pattern searches
- **Parallel processing** using all CPU cores
- **Overflow protection** and **symlink safety** for robust operation
- **Supports all major tech stacks**: Node.js, Python, Java, Go, Rust, Gradle

**Real Performance**: Process 3,000+ files in 5.8ms vs. 15+ seconds with traditional tools!

## 🌳 Interactive Tree View (NEW!)

```bash
./cache-remover-utility --ui ~/Projects
```

**Tree Navigation:**
- `↑↓` or `j/k` - Navigate up/down
- `←→` or `h/l` - Collapse/expand directories
- `Space` - Select projects or directories
- `t` - Toggle between tree and list view
- `a/d` - Select/deselect all
- `c` - Clean selected projects
- `q` - Quit

**Tree Benefits:**
- See directory hierarchy clearly
- Select entire directory trees at once
- Understand project organization
- Aggregate statistics per directory

## 🔧 Advanced Options

```bash
# Scan current directory
./cache-remover-utility

# Limit scan depth
./cache-remover-utility --max-depth 5 ~/Projects

# Use more workers for faster processing
./cache-remover-utility --workers 16 ~/Projects

# Ask before each project cleanup
./cache-remover-utility --interactive ~/Projects

# List all supported project types
./cache-remover-utility --list-types

# Generate config file for customization
./cache-remover-utility --save-config
```

## 📋 What Gets Cleaned

| Technology | Cache Directories | Typical Size Savings |
|------------|-------------------|---------------------|
| **Node.js** | node_modules, dist, build, .next | 100-500 MB per project |
| **Python** | __pycache__, .pytest_cache, venv | 10-100 MB per project |  
| **Java** | target (Maven), build (Gradle) | 50-500 MB per project |
| **Go** | vendor | 10-100 MB per project |
| **Rust** | target | 50-200 MB per project |

## 🆘 Need Help?

- 📖 **Full documentation**: See [USAGE.md](USAGE.md)
- 🐛 **Issues**: [GitHub Issues](https://github.com/abudhahir/projects-cache-clean/issues)
- ❓ **Questions**: Check existing issues or create new one

## ⭐ Quick Tips

**Start with this command to get familiar:**
```bash
./cache-remover-utility --verbose --dry-run ~/Projects
```

**For the best experience, try the new tree view:**
```bash 
./cache-remover-utility --ui ~/Projects
```

**Use the smart run script for convenience:**
```bash 
./run.sh    # Shows interactive menu with all options
```

This will show you exactly what the tool finds and would remove, without actually removing anything!