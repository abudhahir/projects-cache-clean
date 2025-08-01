# ⚡ Quick Start Guide

Get up and running with Cache Remover in under 2 minutes!

## 🚀 Installation

```bash
# Clone and build
git clone https://github.com/abudhahir/projects-cache-clean.git
cd projects-cache-clean
go build -o cache-remover

# Or install Go first if needed
./install-go.sh
```

## 💨 Most Common Usage

```bash
# See what would be cleaned (safe preview)
./cache-remover -dry-run ~/Projects

# Clean everything (with confirmation)
./cache-remover ~/Projects

# Interactive selection with beautiful UI
./cache-remover -ui ~/Projects
```

## 🎬 Visual Examples

### Safe Preview (Always Start Here!)
![Dry Run Demo](screenshots-utility/screenshots/dry-run.gif)

### Basic Cleanup Workflow
![Basic Usage Demo](screenshots-utility/screenshots/basic-usage.gif)

### Interactive TUI Mode
![TUI Interface Demo](screenshots-utility/screenshots/ui-demo.gif)

*These GIFs show the actual terminal output you'll see*

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

## 🎯 Essential Commands

| Command | Purpose |
|---------|---------|
| `./cache-remover ~/Projects` | **Most common** - scan and clean with confirmation |
| `./cache-remover -dry-run ~/Projects` | **Safe preview** - see what would be removed |
| `./cache-remover -ui ~/Projects` | **Interactive** - select projects with TUI |
| `./cache-remover -verbose ~/Projects` | **Detailed** - see all scanning activity |

## 🛡️ Safety First

- **Always start with `-dry-run`** to see what would be removed
- The tool **only removes rebuildable cache** (node_modules, __pycache__, target, etc.)
- **Asks for confirmation** before actual cleanup
- **Skip projects easily** with interactive mode

## ⚡ Performance Features

- **10-100x faster** than traditional tools
- **Smart cache detection** - skips recursing into cache directories
- **Parallel processing** using all CPU cores
- **Supports all major tech stacks**: Node.js, Python, Java, Go, Rust, Gradle

## 🎯 Interactive TUI Mode

```bash
./cache-remover -ui ~/Projects
```

**Controls:**
- `↑↓` - Navigate projects
- `Space` - Select/deselect project
- `a` - Select all projects  
- `c` - Clean selected projects
- `q` - Quit

## 🔧 Quick Options

```bash
# Scan current directory
./cache-remover

# Limit scan depth
./cache-remover -max-depth 5 ~/Projects

# Use more workers for faster processing
./cache-remover -workers 16 ~/Projects

# Ask before each project cleanup
./cache-remover -interactive ~/Projects
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

## ⭐ Quick Tip

**Start with this command to get familiar:**
```bash
./cache-remover -verbose -dry-run ~/Projects
```

This will show you exactly what the tool finds and would remove, without actually removing anything!