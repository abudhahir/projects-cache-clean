#!/usr/bin/env python3
"""
Cache Remover Utility - Python Implementation
A high-performance, parallel cache removal utility for development projects.
"""

import argparse
import concurrent.futures
import json
import os
import shutil
import sys
import time
from pathlib import Path
from typing import Dict, List, NamedTuple, Optional, Set


class CacheConfig(NamedTuple):
    directories: List[str]
    files: List[str]
    extensions: List[str]


class ProjectType(NamedTuple):
    name: str
    indicators: List[str]
    cache_config: CacheConfig


class CacheItem(NamedTuple):
    path: str
    size: int
    type: str


class CleanupStats:
    def __init__(self):
        self.total_projects = 0
        self.total_cache_items = 0
        self.total_size_removed = 0
        self.processing_time = 0


PROJECT_TYPES = [
    ProjectType(
        name="Node.js",
        indicators=["package.json", "yarn.lock", "package-lock.json"],
        cache_config=CacheConfig(
            directories=["node_modules", "dist", "build", ".next", ".nuxt", "coverage"],
            files=[],
            extensions=[]
        )
    ),
    ProjectType(
        name="Python",
        indicators=["requirements.txt", "setup.py", "pyproject.toml", "Pipfile"],
        cache_config=CacheConfig(
            directories=["__pycache__", ".pytest_cache", "dist", "build", ".mypy_cache", ".tox", "venv", ".venv"],
            files=[],
            extensions=[".pyc", ".pyo"]
        )
    ),
    ProjectType(
        name="Java/Maven",
        indicators=["pom.xml"],
        cache_config=CacheConfig(
            directories=["target"],
            files=[],
            extensions=[]
        )
    ),
    ProjectType(
        name="Gradle",
        indicators=["build.gradle", "build.gradle.kts"],
        cache_config=CacheConfig(
            directories=["build", ".gradle"],
            files=[],
            extensions=[]
        )
    ),
    ProjectType(
        name="Go",
        indicators=["go.mod", "go.sum"],
        cache_config=CacheConfig(
            directories=["vendor"],
            files=[],
            extensions=[]
        )
    ),
    ProjectType(
        name="Rust",
        indicators=["Cargo.toml"],
        cache_config=CacheConfig(
            directories=["target"],
            files=[],
            extensions=[]
        )
    ),
]


def format_bytes(bytes_count: int) -> str:
    """Format bytes to human readable format."""
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if bytes_count < 1024.0:
            return f"{bytes_count:.1f} {unit}"
        bytes_count /= 1024.0
    return f"{bytes_count:.1f} PB"


def get_dir_size(path: Path) -> int:
    """Calculate total size of directory."""
    total_size = 0
    try:
        for entry in path.rglob('*'):
            if entry.is_file():
                try:
                    total_size += entry.stat().st_size
                except (OSError, IOError):
                    continue
    except (OSError, IOError):
        pass
    return total_size


def is_project_directory(path: Path) -> bool:
    """Check if directory contains any project indicators."""
    for project_type in PROJECT_TYPES:
        for indicator in project_type.indicators:
            if (path / indicator).exists():
                return True
    return False


def detect_project_type(path: Path) -> Optional[ProjectType]:
    """Detect the type of project in the given directory."""
    for project_type in PROJECT_TYPES:
        for indicator in project_type.indicators:
            if (path / indicator).exists():
                return project_type
    return None


def find_projects(root_dir: Path, max_depth: int, verbose: bool) -> List[Path]:
    """Find all project directories recursively."""
    projects = []
    
    def scan_directory(current_path: Path, depth: int):
        if depth > max_depth:
            return
        
        try:
            if is_project_directory(current_path):
                projects.append(current_path)
                if verbose:
                    print(f"📁 Found project: {current_path}")
            
            # Scan subdirectories
            for entry in current_path.iterdir():
                if entry.is_dir() and not entry.name.startswith('.'):
                    scan_directory(entry, depth + 1)
                    
        except (OSError, IOError, PermissionError):
            if verbose:
                print(f"⚠️  Permission denied: {current_path}")
    
    scan_directory(root_dir, 0)
    return projects


def find_cache_items(project_path: Path, config: CacheConfig) -> List[CacheItem]:
    """Find all cache items in a project directory."""
    items = []
    
    # Check directories
    for dir_name in config.directories:
        dir_path = project_path / dir_name
        if dir_path.exists() and dir_path.is_dir():
            size = get_dir_size(dir_path)
            if size > 0:
                items.append(CacheItem(str(dir_path), size, "directory"))
    
    # Check files
    for file_name in config.files:
        file_path = project_path / file_name
        if file_path.exists() and file_path.is_file():
            try:
                size = file_path.stat().st_size
                items.append(CacheItem(str(file_path), size, "file"))
            except (OSError, IOError):
                pass
    
    # Check extensions
    if config.extensions:
        try:
            for file_path in project_path.rglob('*'):
                if file_path.is_file() and file_path.suffix in config.extensions:
                    try:
                        size = file_path.stat().st_size
                        items.append(CacheItem(str(file_path), size, "file"))
                    except (OSError, IOError):
                        pass
        except (OSError, IOError):
            pass
    
    return items


def remove_cache_items(items: List[CacheItem], verbose: bool) -> tuple[int, int]:
    """Remove cache items and return count and total size removed."""
    removed_items = 0
    removed_size = 0
    
    for item in items:
        try:
            path = Path(item.path)
            if path.exists():
                if item.type == "directory":
                    shutil.rmtree(path)
                else:
                    path.unlink()
                removed_items += 1
                removed_size += item.size
                if verbose:
                    print(f"🗑️  Removed: {item.path} ({format_bytes(item.size)})")
        except (OSError, IOError, PermissionError) as e:
            if verbose:
                print(f"❌ Failed to remove {item.path}: {e}")
    
    return removed_items, removed_size


def process_project(project_path: Path, dry_run: bool, verbose: bool, interactive: bool) -> tuple[int, int]:
    """Process a single project and return (items_removed, size_removed)."""
    project_type = detect_project_type(project_path)
    if not project_type:
        return 0, 0
    
    if verbose:
        print(f"🔍 Processing {project_type.name} project: {project_path}")
    
    cache_items = find_cache_items(project_path, project_type.cache_config)
    if not cache_items:
        if verbose:
            print(f"✅ No cache found in: {project_path}")
        return 0, 0
    
    total_size = sum(item.size for item in cache_items)
    
    print(f"🗂️  {project_path.name} ({project_type.name}): {len(cache_items)} cache items ({format_bytes(total_size)})")
    
    if interactive and not dry_run:
        response = input(f"Remove cache for {project_path}? [y/N]: ").strip().lower()
        if response not in ['y', 'yes']:
            print(f"⏭️  Skipped: {project_path}")
            return 0, 0
    
    if dry_run:
        print(f"🔍 Would remove {len(cache_items)} items ({format_bytes(total_size)}) from: {project_path}")
        for item in cache_items:
            print(f"  - {item.path} ({format_bytes(item.size)})")
        return 0, 0
    else:
        removed_items, removed_size = remove_cache_items(cache_items, verbose)
        if removed_items > 0:
            print(f"✅ Removed {removed_items} items ({format_bytes(removed_size)}) from: {project_path}")
        return removed_items, removed_size


def main():
    parser = argparse.ArgumentParser(description="Cache Remover Utility - Remove rebuildable cache files from projects")
    parser.add_argument("--dir", default=".", help="Root directory to scan for projects")
    parser.add_argument("--dry-run", action="store_true", help="Show what would be removed without actually removing")
    parser.add_argument("--workers", type=int, default=os.cpu_count(), help="Number of worker threads")
    parser.add_argument("--verbose", action="store_true", help="Verbose output")
    parser.add_argument("--max-depth", type=int, default=10, help="Maximum directory depth to scan")
    parser.add_argument("--interactive", action="store_true", help="Ask for confirmation before removing each cache")
    
    args = parser.parse_args()
    
    print("🧹 Cache Remover Utility")
    print(f"Scanning directory: {args.dir}")
    print(f"Workers: {args.workers}")
    if args.dry_run:
        print("🔍 DRY RUN MODE - No files will be removed")
    print()
    
    root_path = Path(args.dir).resolve()
    if not root_path.exists():
        print(f"❌ Directory does not exist: {root_path}")
        sys.exit(1)
    
    start_time = time.time()
    stats = CleanupStats()
    
    # Find all projects
    projects = find_projects(root_path, args.max_depth, args.verbose)
    print(f"Found {len(projects)} projects\n")
    
    if not projects:
        print("No projects found.")
        return
    
    stats.total_projects = len(projects)
    
    # Process projects in parallel
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.workers) as executor:
        futures = [
            executor.submit(process_project, project, args.dry_run, args.verbose, args.interactive)
            for project in projects
        ]
        
        for future in concurrent.futures.as_completed(futures):
            try:
                items_removed, size_removed = future.result()
                stats.total_cache_items += items_removed
                stats.total_size_removed += size_removed
            except Exception as e:
                if args.verbose:
                    print(f"❌ Error processing project: {e}")
            print()  # Add spacing between projects
    
    stats.processing_time = time.time() - start_time
    
    # Print statistics
    print("📊 Cleanup Statistics:")
    print(f"   Projects processed: {stats.total_projects}")
    print(f"   Cache items removed: {stats.total_cache_items}")
    print(f"   Total space reclaimed: {format_bytes(stats.total_size_removed)}")
    print(f"   Processing time: {stats.processing_time:.2f}s")
    if stats.processing_time > 0:
        speed_mb_s = (stats.total_size_removed / (1024 * 1024)) / stats.processing_time
        print(f"   Average speed: {speed_mb_s:.2f} MB/s")


if __name__ == "__main__":
    main()