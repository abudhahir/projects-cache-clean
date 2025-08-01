# 📖 Comprehensive Usage Guide

Complete documentation for all features and advanced usage of Cache Remover Utility.

## 📋 Table of Contents

- [Command Line Options](#-command-line-options)
- [Usage Modes](#-usage-modes)  
- [Interactive TUI Guide](#-interactive-tui-guide)
- [Advanced Examples](#-advanced-examples)
- [Project Type Detection](#-project-type-detection)
- [Performance Optimizations](#-performance-optimizations)
- [Safety Features](#️-safety-features)
- [Troubleshooting](#-troubleshooting)

## 🔧 Command Line Options

### Core Options
| Flag | Default | Description |
|------|---------|-------------|
| `directory` | `.` | **Root directory to scan** (positional argument) |
| `-dir` | `.` | Alternative flag for root directory |
| `-dry-run` | `false` | Show what would be removed without removing |
| `-verbose` | `false` | Verbose output with detailed logging |

### Interface Options  
| Flag | Default | Description |
|------|---------|-------------|
| `-ui` | `false` | **Launch interactive terminal UI** |
| `-interactive` | `false` | Ask for confirmation before removing each cache |

### Performance Options
| Flag | Default | Description |
|------|---------|-------------|
| `-workers` | CPU cores | Number of worker goroutines |
| `-max-depth` | `10` | Maximum directory depth to scan |

## 🎬 Visual Examples

### Dry Run Mode - Safe Preview
![Dry Run Demo](screenshots-utility/screenshots/dry-run.gif)

### Verbose Mode with Optimization Messages
![Verbose Mode Demo](screenshots-utility/screenshots/verbose.gif)

### Interactive Per-Project Confirmation
![Interactive Mode Demo](screenshots-utility/screenshots/interactive.gif)

### Performance Optimization in Action
![Performance Demo](screenshots-utility/screenshots/performance.gif)

*These GIFs show the actual terminal output you'll see when using different modes*

## 🎛️ Usage Modes

### 1. 📊 Preview Mode (Recommended First Step)
```bash
# See what would be cleaned without making changes
./cache-remover -dry-run ~/Projects

# Verbose preview with detailed scanning info
./cache-remover -verbose -dry-run ~/Projects
```

**Output:**
```
🧹 Cache Remover Utility
Scanning directory: /Users/dev/Projects
Workers: 8
🔍 DRY RUN MODE - No files will be removed

📁 Found project: /Users/dev/Projects/my-app
⏭️  Skipping cache directory: /Users/dev/Projects/my-app/node_modules
Found 1 projects

🗂️  my-app (Node.js): 2 cache items (145.2 MB)  
🔍 Would remove 2 items (145.2 MB) from: /Users/dev/Projects/my-app
  - /Users/dev/Projects/my-app/node_modules (142.1 MB)
  - /Users/dev/Projects/my-app/build (3.1 MB)

📊 Cleanup Statistics:
   Projects processed: 1
   Cache items removed: 2  
   Total space reclaimed: 145.2 MB
   Processing time: 0.8s
   Average speed: 181.5 MB/s
```

### 2. ⚡ Automatic Mode (Quick Cleanup)
```bash
# Clean all projects with single confirmation
./cache-remover ~/Projects

# Silent mode (minimal output)
./cache-remover ~/Projects 2>/dev/null
```

**Output:**
```
🧹 Cache Remover Utility
Scanning directory: /Users/dev/Projects
Workers: 8

Found 3 projects

🗂️  my-react-app (Node.js): 2 cache items (245.2 MB)
🗂️  python-api (Python): 3 cache items (12.8 MB)  
🗂️  java-service (Java/Maven): 1 cache items (156.4 MB)

This will remove cache files totaling 414.4 MB from 3 projects.
Continue? [y/N]: y

✅ Removed 2 items (245.2 MB) from: /Users/dev/Projects/my-react-app
✅ Removed 3 items (12.8 MB) from: /Users/dev/Projects/python-api  
✅ Removed 1 items (156.4 MB) from: /Users/dev/Projects/java-service

📊 Cleanup Statistics:
   Projects processed: 3
   Cache items removed: 6
   Total space reclaimed: 414.4 MB
   Processing time: 3.1s
   Average speed: 133.68 MB/s
```

### 3. 🤔 Interactive Mode (Per-Project Confirmation)
```bash
# Ask before cleaning each individual project
./cache-remover -interactive ~/Projects
```

**Output:**
```
🗂️  my-react-app (Node.js): 2 cache items (245.2 MB)
Remove cache for /Users/dev/Projects/my-react-app? [y/N]: y
✅ Removed 2 items (245.2 MB) from: /Users/dev/Projects/my-react-app

🗂️  python-api (Python): 3 cache items (12.8 MB)  
Remove cache for /Users/dev/Projects/python-api? [y/N]: n
⏭️  Skipped: /Users/dev/Projects/python-api

🗂️  java-service (Java/Maven): 1 cache items (156.4 MB)
Remove cache for /Users/dev/Projects/java-service? [y/N]: y
✅ Removed 1 items (156.4 MB) from: /Users/dev/Projects/java-service
```

## 🖥️ Interactive TUI Guide

### TUI Interface Demo
![TUI Interface Demo](screenshots-utility/screenshots/ui-demo.gif)

### Launching TUI
```bash
./cache-remover -ui ~/Projects
```

### TUI Interface Layout
```
🚀 Launching Interactive TUI Cache Remover...

┌─ 🧹 Cache Remover - Project Selection ─────────────────────────────────┐
│                                                                         │
│  📁 my-react-app (Node.js)                          [✓] 245.2 MB       │
│     └─ 2 cache items: node_modules, build                              │
│                                                                         │
│  📁 python-api (Python)                             [ ] 12.8 MB        │
│     └─ 3 cache items: __pycache__, .pytest_cache, build               │
│                                                                         │
│  📁 java-service (Java/Maven)                       [✓] 156.4 MB       │
│     └─ 1 cache items: target                                           │
│                                                                         │
│  📊 Statistics: 3 projects, 2 selected, 401.6 MB reclaimable          │
│                                                                         │
│  🔧 Controls: ↑↓=navigate, space=select, a=select all, c=clean, q=quit │
└─────────────────────────────────────────────────────────────────────────┘
```

### TUI Keyboard Shortcuts
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up in project list |
| `↓` / `j` | Move down in project list |
| `Space` / `Enter` | Toggle project selection |
| `a` | Select all projects |
| `d` | Deselect all projects |
| `c` | Clean selected projects |
| `r` | Refresh project list |
| `v` | View detailed project information |
| `?` | Show help/shortcuts |
| `q` / `Esc` | Quit application |

### TUI Workflow
1. **Navigate** projects with arrow keys
2. **Select** projects you want to clean with `Space`
3. **Review** the statistics at the bottom
4. **Clean** selected projects with `c`
5. **Confirm** the cleanup operation
6. **Watch** progress indicators during cleanup

## 💡 Advanced Examples

### Performance Tuning
```bash
# Use maximum workers for large directories
./cache-remover -workers 32 ~/Projects

# Limit depth for faster scanning of shallow structures
./cache-remover -max-depth 3 ~/Projects

# Combine for optimal performance
./cache-remover -workers 16 -max-depth 5 -verbose ~/Projects
```

### Targeted Scanning
```bash
# Scan only specific subdirectories
./cache-remover ~/Projects/frontend-projects
./cache-remover ~/Projects/backend-services

# Scan current directory
./cache-remover .

# Scan parent directory
./cache-remover ..
```

### Scripting Integration
```bash
#!/bin/bash
# Automated cleanup script

echo "🧹 Starting automated cache cleanup..."

# Preview what would be cleaned
./cache-remover -dry-run ~/Projects > cleanup-preview.log

# Extract total space from preview
SPACE=$(grep "Total space reclaimed" cleanup-preview.log | awk '{print $4}')

echo "Found $SPACE of cache files to clean"

# Auto-confirm cleanup if space > 100MB
if [[ "$SPACE" =~ "GB" ]] || [[ "$SPACE" =~ "[0-9]{3,}" ]]; then
    echo "Proceeding with cleanup..."
    echo "y" | ./cache-remover ~/Projects
else
    echo "Less than 100MB found, skipping cleanup"
fi
```

### Output Redirection
```bash
# Save detailed log
./cache-remover -verbose ~/Projects > cache-cleanup.log 2>&1

# Only show errors
./cache-remover ~/Projects 2>error.log >/dev/null

# Quiet mode (only final statistics)
./cache-remover ~/Projects | grep "📊"
```

## 🔍 Project Type Detection

### Detection Logic
The tool identifies projects by looking for specific indicator files:

```bash
# Node.js Projects - looks for:
package.json, yarn.lock, package-lock.json

# Python Projects - looks for:  
requirements.txt, setup.py, pyproject.toml, Pipfile

# Java Projects - looks for:
pom.xml (Maven), build.gradle, build.gradle.kts (Gradle)

# Go Projects - looks for:
go.mod, go.sum

# Rust Projects - looks for:
Cargo.toml
```

### Cache Patterns by Technology

#### Node.js Cache Cleanup
```
📁 my-react-app/
├── node_modules/          ← Removed (packages)
├── dist/                  ← Removed (build output)  
├── build/                 ← Removed (build output)
├── .next/                 ← Removed (Next.js cache)
├── coverage/              ← Removed (test coverage)
├── src/                   ← Kept (source code)
└── package.json           ← Kept (configuration)
```

#### Python Cache Cleanup  
```
📁 python-api/
├── __pycache__/           ← Removed (compiled Python)
├── .pytest_cache/         ← Removed (test cache)
├── .mypy_cache/           ← Removed (type checker cache)
├── dist/                  ← Removed (distribution packages)
├── venv/                  ← Removed (virtual environment)
├── src/                   ← Kept (source code)
└── requirements.txt       ← Kept (dependencies)
```

#### Java Cache Cleanup
```  
📁 java-service/
├── target/                ← Removed (Maven build output)
├── build/                 ← Removed (Gradle build output)
├── .gradle/               ← Removed (Gradle cache)
├── src/                   ← Kept (source code)
└── pom.xml                ← Kept (Maven configuration)
```

## ⚡ Performance Optimizations

### Smart Cache Directory Skipping
The tool implements intelligent scanning optimizations:

**Before Optimization:**
```
📁 node_modules/
├── 📄 scanning package1/index.js
├── 📄 scanning package1/lib/util.js  
├── 📄 scanning package1/lib/helper.js
├── 📄 scanning package2/index.js
└── ... (scanning 50,000+ files) ← SLOW!
```

**After Optimization:**
```
📁 node_modules/           ← Treat as single unit ← FAST!
⏭️  Skipping cache directory: /path/to/node_modules
```

### Performance Comparison
| Scenario | Before | After | Improvement |
|----------|--------|-------|-------------|
| Large node_modules (50k files) | 28.3s | 0.1s | **283x faster** |
| Maven target (15k files) | 8.7s | 0.05s | **174x faster** |
| Python __pycache__ (1k files) | 2.1s | 0.01s | **210x faster** |

### Memory Efficiency
```bash
# Memory usage comparison
Before: ~500MB RAM (loading full directory tree)
After:  ~30MB RAM (streaming traversal)
```

## 🛡️ Safety Features

### 1. Dry Run Protection
```bash
# Always safe to run - shows preview without changes
./cache-remover -dry-run ~/Projects
```

### 2. Project Type Validation
```bash
# Only removes caches from recognized project types
# Ignores random directories without project indicators
```

### 3. Confirmation Prompts
```bash
# Global confirmation
"This will remove cache files totaling 414.4 MB from 3 projects. Continue? [y/N]:"

# Per-project confirmation (interactive mode)
"Remove cache for /path/to/project? [y/N]:"
```

### 4. Error Handling
```bash
# Graceful handling of permission issues
❌ Failed to remove /protected/cache: permission denied
✅ Continuing with other projects...
```

### 5. Depth Limiting
```bash
# Prevents infinite recursion with symlinks
-max-depth 10  # Default limit prevents runaway scanning
```

## 🔧 Troubleshooting

### Common Issues

#### "Permission denied" errors
```bash
# Run with appropriate permissions
sudo ./cache-remover ~/Projects

# Or fix directory permissions
chmod -R u+w ~/Projects
```

#### "No projects found"
```bash
# Check if directory contains recognizable project types
ls ~/Projects/*/package.json  # Node.js
ls ~/Projects/*/requirements.txt  # Python  
ls ~/Projects/*/pom.xml  # Java/Maven

# Try verbose mode to see scanning details
./cache-remover -verbose -dry-run ~/Projects
```

#### TUI not working
```bash
# Ensure terminal supports TUI
export TERM=xterm-256color

# Try regular mode instead
./cache-remover ~/Projects
```

#### Slow performance
```bash
# Reduce worker count for slower systems
./cache-remover -workers 4 ~/Projects

# Limit scanning depth
./cache-remover -max-depth 5 ~/Projects

# Check if optimization is working (should see "Skipping" messages)
./cache-remover -verbose ~/Projects | grep "Skipping"
```

### Performance Debugging
```bash
# Check scanning efficiency
./cache-remover -verbose -dry-run ~/Projects | grep "⏭️  Skipping"

# Should see messages like:
⏭️  Skipping cache directory: /path/to/node_modules
⏭️  Skipping cache directory: /path/to/__pycache__
⏭️  Skipping cache directory: /path/to/target
```

### Logging and Debugging
```bash
# Full verbose output with timing
time ./cache-remover -verbose ~/Projects

# Capture all output for analysis
./cache-remover -verbose ~/Projects > full-log.txt 2>&1

# Extract just statistics
./cache-remover ~/Projects | grep "📊" -A 10
```

## 📊 Output Format Reference

### Statistics Block
```
📊 Cleanup Statistics:
   Projects processed: 3        ← Number of projects scanned
   Cache items removed: 6       ← Individual cache files/directories removed  
   Total space reclaimed: 414.4 MB  ← Actual disk space freed
   Processing time: 3.1s        ← Total execution time
   Average speed: 133.68 MB/s   ← Removal throughput
```

### Progress Indicators
```
📁 Found project: /path          ← Project discovery
⏭️  Skipping cache directory     ← Optimization in action
🔍 Processing ProjectType        ← Project being analyzed
🗂️  ProjectName (Type): N items ← Cache summary
✅ Removed N items (Size)       ← Successful cleanup
⏭️  Skipped: /path              ← User-skipped project
❌ Failed to remove             ← Error (with details)
```

## 🚀 Integration Examples

### CI/CD Pipeline
```yaml
# .github/workflows/cleanup.yml
name: Weekly Cache Cleanup
on:
  schedule:
    - cron: '0 2 * * 0'  # Sunday 2AM
    
jobs:
  cleanup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Build cache remover
        run: go build -o cache-remover
      - name: Clean caches
        run: echo "y" | ./cache-remover .
```

### Docker Integration
```dockerfile
FROM golang:1.21-alpine AS builder
COPY . .
RUN go build -o cache-remover

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /cache-remover /usr/local/bin/
ENTRYPOINT ["cache-remover"]
```

### Makefile Integration
```makefile
.PHONY: clean-cache
clean-cache:
	@echo "🧹 Cleaning project caches..."
	@./cache-remover -dry-run .
	@read -p "Proceed with cleanup? [y/N]: " confirm && \
	 [[ $$confirm == [yY] ]] && ./cache-remover . || echo "Cancelled"
```

---

For more help, check the [QUICKSTART.md](QUICKSTART.md) or create an issue on [GitHub](https://github.com/abudhahir/projects-cache-clean/issues).