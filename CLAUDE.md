# Cache Remover Utility - Project Memory

## Project Overview
Go-based command-line utility for cleaning cache directories across multiple technology stacks.

## Key Technical Architecture
- **Language**: Go
- **Build**: `go build -o cache-remover`
- **Entry Point**: `main.go` with positional directory argument support
- **Core Logic**: `interactive.go` with TUI interface using Bubble Tea
- **Optimization**: Cache directory skipping for 10-100x performance improvement

## Important Project Rules
- **Screenshot Utility**: Located on `feature/screenshot-utility` branch - NEVER merge to main
- **Branch Strategy**: Feature branches for utilities, main branch for core functionality only
- **Dependencies**: VHS and ttyd for screenshot generation (separate branch only)

## Core Features Implemented
1. **Positional Arguments**: `./cache-remover /path/to/directory`
2. **Dry Run Mode**: `-dry-run` flag for safe preview
3. **Verbose Mode**: `-verbose` flag with detailed scanning logs
4. **Cache Directory Optimization**: Skips recursive scanning of known cache dirs
5. **Statistics Collection**: Shows reclaimable space even in dry-run mode
6. **Interactive TUI**: Bubble Tea interface for project selection

## Cache Directories Detected
- `node_modules/` (Node.js)
- `__pycache__/`, `.pytest_cache/`, `build/`, `dist/` (Python)
- `target/` (Java/Maven, Rust)
- `build/` (General build artifacts)

## Performance Optimizations
- Skip descending into cache directories (they're removed as units)
- Early termination when cache directory is identified
- Efficient file system traversal

## Documentation Structure
- **README.md**: Main project overview
- **QUICKSTART.md**: 2-minute getting started guide
- **USAGE.md**: Comprehensive command-line reference
- **Screenshots**: Generated via VHS on separate branch (feature/screenshot-utility)

## Development History
1. Migrated from Python implementation to Go-only
2. Removed webui components for CLI focus
3. Fixed positional argument handling
4. Implemented performance optimizations
5. Created comprehensive documentation
6. Built screenshot generation utility (separate branch)

## Build and Test Commands
```bash
go build -o cache-remover
./cache-remover -dry-run demo-projects
./cache-remover -verbose demo-projects
```

## Current Status
- Core functionality: Complete
- Documentation: Complete with screenshot placeholders
- Screenshot utility: Available on feature/screenshot-utility branch (not for main merge)
- Ready for additional feature development on main branch