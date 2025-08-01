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

## URGENT: Immediate Fix Plan Required
**CRITICAL DISCOVERY**: Codebase evolution analysis revealed that recent "improvements" to progress implementation actually DEGRADED code quality:

### Current State Issues (Must Fix Immediately):
- **Security Vulnerabilities**: Introduced race conditions and memory leaks that didn't exist originally
- **Architecture Violations**: Added global state anti-patterns violating Bubble Tea framework
- **Resource Management**: Created timer leaks and resource management failures
- **Thread Safety**: Original was perfect, current has race conditions

### Original Implementation Was BETTER:
The first commit had a clean, safe, simple progress implementation that worked perfectly. Our "real-time progress" attempt introduced critical problems without substantial benefit.

### Immediate Action Plan:
**Follow IMMEDIATE_FIX.md step by step:**
1. **Issue #1**: Revert problematic progress implementation to original clean version
2. **Issue #2**: Preserve infrastructure improvements (config, docs, tests)  
3. **Issue #3**: Restore original simple progress display
4. **Issue #4**: Add regression prevention measures
5. **Issue #5**: Update documentation to reflect revert
6. **Issue #6**: Validate clean state after revert

### What to KEEP (Excellent improvements):
- ✅ Configuration system (config.go, JSON configuration)
- ✅ Documentation (README.md, QUICKSTART.md, USAGE.md)
- ✅ Enhanced test coverage (290% increase)
- ✅ CLI enhancements (--list-types, --save-config)

### What to REVERT (Introduced problems):
- ❌ Global cleanup state (`var globalCleanupState *cleanupState`)
- ❌ Timer management (`tea.Every()`)
- ❌ Complex progress tracking with race conditions
- ❌ Resource leak prone implementations

## Core Features Implemented (Known Good)
1. **Positional Arguments**: `./cache-remover /path/to/directory`
2. **Dry Run Mode**: `-dry-run` flag for safe preview
3. **Verbose Mode**: `-verbose` flag with detailed scanning logs
4. **Cache Directory Optimization**: Skips recursive scanning of known cache dirs
5. **Statistics Collection**: Shows reclaimable space even in dry-run mode
6. **Interactive TUI**: Bubble Tea interface for project selection
7. **Configuration System**: JSON-based project type configuration
8. **Extended Project Support**: 9 project types (Node.js, Python, Java, etc.)

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
- **CODE_REVIEW_REPORT.md**: Initial high-priority improvements review
- **PROGRESS_IMPLEMENTATION_REVIEW.md**: Critical issues in progress implementation
- **CODEBASE_EVOLUTION_ANALYSIS.md**: Comparative study showing quality degradation
- **IMMEDIATE_FIX.md**: Step-by-step plan to fix critical issues
- **Screenshots**: Generated via VHS on separate branch (feature/screenshot-utility)

## Development History & Lessons Learned
1. Migrated from Python implementation to Go-only ✅
2. Removed webui components for CLI focus ✅
3. Fixed positional argument handling ✅
4. Implemented performance optimizations ✅
5. Created comprehensive documentation ✅
6. Built screenshot generation utility (separate branch) ✅
7. **Added configuration system** ✅ **EXCELLENT**
8. **Enhanced test coverage** ✅ **EXCELLENT**
9. **Attempted progress improvements** ❌ **DEGRADED QUALITY - NEEDS IMMEDIATE REVERT**

## Build and Test Commands
```bash
go build -o cache-remover
./cache-remover -dry-run demo-projects
./cache-remover -verbose demo-projects
go test -v  # Should pass after revert
```

## Critical Lesson Learned
**"Perfect is the enemy of good"** - The original simple implementation was working perfectly. Attempts to add "real-time progress" introduced significant security and architectural problems without substantial benefit.

**Key Principle**: Always make the smallest possible change that works. Don't fix what isn't broken.

## Current Status
- Core functionality: **NEEDS IMMEDIATE REVERT** (progress implementation)
- Infrastructure: **EXCELLENT** (config, docs, tests)
- Overall: **MIXED** - great infrastructure improvements marred by core functionality degradation
- Next action: **Execute IMMEDIATE_FIX.md plan step by step**

## Memory Note
**ALWAYS FOLLOW IMMEDIATE_FIX.md** - The step-by-step plan to restore code quality while preserving infrastructure improvements. This is the highest priority task for maintaining code quality and user data safety.