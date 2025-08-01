package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type CacheConfig struct {
	Directories []string `json:"directories"`
	Files       []string `json:"files"`
	Extensions  []string `json:"extensions"`
}

type ProjectType struct {
	Name        string      `json:"name"`
	Indicators  []string    `json:"indicators"`
	CacheConfig CacheConfig `json:"cache_config"`
}

type CacheItem struct {
	Path string
	Size int64
	Type string
}

type CleanupStats struct {
	TotalProjects    int
	TotalCacheItems  int
	TotalSizeRemoved int64
	ProcessingTime   time.Duration
	mu               sync.Mutex
}

func (s *CleanupStats) Add(items int, size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalCacheItems += items
	s.TotalSizeRemoved += size
}

func (s *CleanupStats) IncrementProjects() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalProjects++
}

func main() {
	// Load configuration first
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	var (
		rootDir     = flag.String("dir", ".", "Root directory to scan for projects")
		dryRun      = flag.Bool("dry-run", false, "Show what would be removed without actually removing")
		workers     = flag.Int("workers", config.Settings.DefaultWorkers, "Number of worker goroutines")
		verbose     = flag.Bool("verbose", false, "Verbose output")
		maxDepth    = flag.Int("max-depth", config.Settings.MaxDepth, "Maximum directory depth to scan")
		interactive = flag.Bool("interactive", false, "Ask for confirmation before removing each cache")
		ui          = flag.Bool("ui", false, "Launch interactive TUI mode")
		saveConfig  = flag.Bool("save-config", false, "Save default configuration to current directory")
		listTypes   = flag.Bool("list-types", false, "List all supported project types")
	)
	flag.Parse()

	// Handle special flags
	if *saveConfig {
		if err := saveDefaultConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Default configuration saved to cache-remover-config.json")
		return
	}

	if *listTypes {
		listProjectTypes(config)
		return
	}

	// Handle positional argument for directory
	if len(flag.Args()) > 0 {
		*rootDir = flag.Args()[0]
	}

	// Launch interactive TUI if requested
	if *ui {
		fmt.Println("🚀 Launching Interactive TUI Cache Remover...")
		if err := runInteractiveUI(*rootDir); err != nil {
			fmt.Printf("Error running interactive UI: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("🧹 Cache Remover Utility\n")
	fmt.Printf("Scanning directory: %s\n", *rootDir)
	fmt.Printf("Workers: %d\n", *workers)
	if *dryRun {
		fmt.Printf("🔍 DRY RUN MODE - No files will be removed\n")
	}
	
	// Display supported project types for transparency
	fmt.Printf("🔧 Supported project types: ")
	var typeNames []string
	for _, pt := range config.ProjectTypes {
		typeNames = append(typeNames, pt.Name)
	}
	fmt.Printf("%s\n", strings.Join(typeNames, ", "))
	
	fmt.Printf("💡 Tip: Use --ui for terminal interface or --web for browser interface\n")
	fmt.Println()

	startTime := time.Now()
	stats := &CleanupStats{}

	projects := findProjects(*rootDir, *maxDepth, *verbose)
	fmt.Printf("Found %d projects\n\n", len(projects))

	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	processProjects(projects, *workers, *dryRun, *verbose, *interactive, stats)

	stats.ProcessingTime = time.Since(startTime)
	printStats(stats)
}

func isCacheDirectory(dirName string) bool {
	config, _ := loadConfig()
	// Check if this directory name matches any known cache directory patterns
	for _, projectType := range config.ProjectTypes {
		for _, cacheDir := range projectType.CacheConfig.Directories {
			if dirName == cacheDir {
				return true
			}
		}
	}
	return false
}

func findProjects(rootDir string, maxDepth int, verbose bool) []string {
	var projects []string
	var mu sync.Mutex

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if verbose {
				fmt.Printf("⚠️  Warning: Cannot access %s: %v\n", path, err)
			}
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		depth := strings.Count(strings.TrimPrefix(path, rootDir), string(os.PathSeparator))
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Skip descending into cache directories - they're meant to be removed as units
		if isCacheDirectory(info.Name()) {
			if verbose {
				fmt.Printf("⏭️  Skipping cache directory: %s\n", path)
			}
			return filepath.SkipDir
		}

		if isProjectDirectory(path) {
			mu.Lock()
			projects = append(projects, path)
			mu.Unlock()
			if verbose {
				fmt.Printf("📁 Found project: %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning directories: %v\n", err)
	}

	return projects
}

func isProjectDirectory(dir string) bool {
	config, _ := loadConfig()
	for _, projectType := range config.ProjectTypes {
		for _, indicator := range projectType.Indicators {
			indicatorPath := filepath.Join(dir, indicator)
			if _, err := os.Stat(indicatorPath); err == nil {
				return true
			}
		}
	}
	return false
}

func processProjects(projects []string, workers int, dryRun, verbose, interactive bool, stats *CleanupStats) {
	projectChan := make(chan string, len(projects))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for project := range projectChan {
				processProject(project, dryRun, verbose, interactive, stats)
			}
		}()
	}

	for _, project := range projects {
		projectChan <- project
	}
	close(projectChan)

	wg.Wait()
}

func processProject(projectPath string, dryRun, verbose, interactive bool, stats *CleanupStats) {
	projectType := detectProjectType(projectPath)
	if projectType == nil {
		return
	}

	stats.IncrementProjects()

	if verbose {
		fmt.Printf("🔍 Processing %s project: %s\n", projectType.Name, projectPath)
	}

	cacheItems := findCacheItems(projectPath, projectType.CacheConfig)
	if len(cacheItems) == 0 {
		if verbose {
			fmt.Printf("✅ No cache found in: %s\n", projectPath)
		}
		return
	}

	totalSize := int64(0)
	for _, item := range cacheItems {
		totalSize += item.Size
	}

	fmt.Printf("🗂️  %s (%s): %d cache items (%s)\n",
		filepath.Base(projectPath),
		projectType.Name,
		len(cacheItems),
		formatBytes(totalSize))

	if interactive && !dryRun {
		fmt.Printf("Remove cache for %s? [y/N]: ", projectPath)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Printf("⏭️  Skipped: %s\n", projectPath)
			return
		}
	}

	if dryRun {
		fmt.Printf("🔍 Would remove %d items (%s) from: %s\n",
			len(cacheItems), formatBytes(totalSize), projectPath)
		for _, item := range cacheItems {
			fmt.Printf("  - %s (%s)\n", item.Path, formatBytes(item.Size))
		}
		// Add to stats even in dry-run mode to show potential savings
		stats.Add(len(cacheItems), totalSize)
	} else {
		removedItems, removedSize := removeCacheItems(cacheItems, verbose)
		if removedItems > 0 {
			fmt.Printf("✅ Removed %d items (%s) from: %s\n",
				removedItems, formatBytes(removedSize), projectPath)
		}
		stats.Add(removedItems, removedSize)
	}
	fmt.Println()
}

func detectProjectType(projectPath string) *ProjectType {
	config, _ := loadConfig()
	for _, projectType := range config.ProjectTypes {
		for _, indicator := range projectType.Indicators {
			indicatorPath := filepath.Join(projectPath, indicator)
			if _, err := os.Stat(indicatorPath); err == nil {
				return &projectType
			}
		}
	}
	return nil
}

func findCacheItems(projectPath string, config CacheConfig) []CacheItem {
	var items []CacheItem
	processedPaths := make(map[string]bool) // Track paths to avoid double-counting

	// First: Collect cache directories (search recursively for cache directory names)
	for _, dir := range config.Directories {
		// Check root level first (most common case)
		dirPath := filepath.Join(projectPath, dir)
		if size := getDirSize(dirPath); size > 0 {
			items = append(items, CacheItem{
				Path: dirPath,
				Size: size,
				Type: "directory",
			})
			processedPaths[dirPath] = true
		}
		
		// Then search recursively for nested cache directories
		filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() {
				return nil
			}
			
			// Skip if already processed
			if processedPaths[path] {
				return filepath.SkipDir
			}
			
			// Check if this directory matches a cache directory name
			if info.Name() == dir {
				if size := getDirSize(path); size > 0 {
					items = append(items, CacheItem{
						Path: path,
						Size: size,
						Type: "directory",
					})
					processedPaths[path] = true
					return filepath.SkipDir // Don't traverse into this cache directory
				}
			}
			
			return nil
		})
	}

	// Second: Collect individual cache files
	for _, file := range config.Files {
		filePath := filepath.Join(projectPath, file)
		if !processedPaths[filePath] {
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				items = append(items, CacheItem{
					Path: filePath,
					Size: info.Size(),
					Type: "file",
				})
			}
		}
	}

	// Third: Collect files by extension (but skip those already inside processed cache directories)
	if len(config.Extensions) > 0 {
		filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Skip if this file is inside an already-processed cache directory
			for processedDir := range processedPaths {
				if strings.HasPrefix(path, processedDir+string(filepath.Separator)) {
					return nil // Skip files inside cache directories we already counted
				}
			}

			// Process files with matching extensions
			if !info.IsDir() {
				for _, ext := range config.Extensions {
					if strings.HasSuffix(info.Name(), ext) {
						items = append(items, CacheItem{
							Path: path,
							Size: info.Size(),
							Type: "file",
						})
						break
					}
				}
			}
			return nil
		})
	}

	return items
}

func getDirSize(dirPath string) int64 {
	// Always use full recursive scanning for accurate size calculations
	// Performance optimization happens during project discovery, not size calculation

	var size int64
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files we can't access but continue walking
			return nil
		}
		if !info.IsDir() {
			// Check for potential overflow
			if size > 0 && info.Size() > 0 && size > (9223372036854775807-info.Size()) {
				return filepath.SkipDir // Stop walking to prevent overflow
			}
			size += info.Size()
		}
		return nil
	})
	return size
}

// getOptimizedCacheDirSize calculates directory size without walking through all contents
// This provides significant performance improvement for large cache directories like node_modules
func getOptimizedCacheDirSize(dirPath string) int64 {
	// Check if directory exists and is accessible
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return 0
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// If we can't read the directory (permissions, etc.), fall back to standard method
		return getDirSizeFallback(dirPath)
	}

	var totalSize int64
	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())

		// Skip symlinks to avoid infinite loops and potential issues
		if entry.Type()&os.ModeSymlink != 0 {
			continue
		}

		if entry.IsDir() {
			// For subdirectories, get their size recursively
			// but limit depth to avoid performance issues
			totalSize += getDirectorySizeWithLimit(entryPath, 3)
		} else {
			// For files, get their size directly
			if info, err := entry.Info(); err == nil {
				// Ensure we don't overflow int64
				if totalSize > 0 && info.Size() > 0 && totalSize > (9223372036854775807-info.Size()) {
					// Handle potential overflow
					return totalSize
				}
				totalSize += info.Size()
			}
		}
	}

	return totalSize
}

// getDirectorySizeWithLimit calculates directory size with a maximum depth limit
func getDirectorySizeWithLimit(dirPath string, maxDepth int) int64 {
	if maxDepth <= 0 {
		return 0
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		// Skip directories we can't access
		return 0
	}

	var totalSize int64
	for _, entry := range entries {
		// Skip symlinks to prevent infinite loops
		if entry.Type()&os.ModeSymlink != 0 {
			continue
		}

		if !entry.IsDir() {
			if info, err := entry.Info(); err == nil {
				// Check for potential overflow
				if totalSize > 0 && info.Size() > 0 && totalSize > (9223372036854775807-info.Size()) {
					return totalSize
				}
				totalSize += info.Size()
			}
		} else if maxDepth > 1 {
			// Recurse into subdirectories with reduced depth
			subPath := filepath.Join(dirPath, entry.Name())
			subSize := getDirectorySizeWithLimit(subPath, maxDepth-1)

			// Check for potential overflow
			if totalSize > 0 && subSize > 0 && totalSize > (9223372036854775807-subSize) {
				return totalSize
			}
			totalSize += subSize
		}
	}

	return totalSize
}

// getDirSizeFallback is the original directory size calculation method
// Used as fallback when optimized method fails
func getDirSizeFallback(dirPath string) int64 {
	var size int64
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip files we can't access but continue walking
			return nil
		}
		if !info.IsDir() {
			// Check for potential overflow
			if size > 0 && info.Size() > 0 && size > (9223372036854775807-info.Size()) {
				return filepath.SkipDir // Stop walking to prevent overflow
			}
			size += info.Size()
		}
		return nil
	})
	return size
}

func removeCacheItems(items []CacheItem, verbose bool) (int, int64) {
	removedItems := 0
	removedSize := int64(0)

	for _, item := range items {
		if err := forceRemoveCacheDirectory(item.Path, verbose); err != nil {
			// Always log removal failures, not just in verbose mode
			fmt.Printf("❌ Failed to remove %s: %v\n", item.Path, err)
		} else {
			removedItems++
			removedSize += item.Size
			if verbose {
				fmt.Printf("🗑️  Removed: %s (%s)\n", item.Path, formatBytes(item.Size))
			}
		}
	}

	return removedItems, removedSize
}

// forceRemoveCacheDirectory aggressively removes cache directories with multiple strategies
func forceRemoveCacheDirectory(path string, verbose bool) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Already gone, consider it success
	}

	// Strategy 1: Try standard RemoveAll first (fastest when it works)
	if err := os.RemoveAll(path); err == nil {
		return nil
	}

	// Strategy 2: Try to fix permissions and retry
	if err := makeWritableRecursive(path); err == nil {
		if err := os.RemoveAll(path); err == nil {
			return nil
		}
	}

	// Strategy 3: Manual recursive removal with permission fixing
	if err := forceRemoveRecursive(path, verbose); err == nil {
		return nil
	}

	// Strategy 4: Last resort - try system command (Unix-like systems only)
	if err := forceRemoveWithSystemCommand(path); err == nil {
		return nil
	}

	// If all strategies fail, return the error
	return fmt.Errorf("all removal strategies failed for cache directory: %s", path)
}

// makeWritableRecursive makes all files and directories writable
func makeWritableRecursive(path string) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		
		// Make file/directory writable
		return os.Chmod(filePath, info.Mode()|0200) // Add write permission
	})
}

// forceRemoveRecursive manually removes files and directories with permission fixing
func forceRemoveRecursive(path string, verbose bool) error {
	// First pass: fix all permissions
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		// Make writable and executable (for directories)
		newMode := info.Mode() | 0200 // Add write permission
		if info.IsDir() {
			newMode |= 0100 // Add execute permission for directories
		}
		os.Chmod(filePath, newMode)
		return nil
	})

	// Second pass: remove everything (start from deepest level)
	var allPaths []string
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err == nil {
			allPaths = append(allPaths, filePath)
		}
		return nil
	})

	// Sort paths by depth (deepest first) to avoid "directory not empty" errors
	// Reverse the slice so we process deepest paths first
	for i := len(allPaths) - 1; i >= 0; i-- {
		filePath := allPaths[i]
		if info, err := os.Stat(filePath); err == nil {
			if info.IsDir() {
				if err := os.Remove(filePath); err != nil && verbose {
					fmt.Printf("⚠️  Warning: Cannot remove directory %s: %v\n", filePath, err)
				}
			} else {
				if err := os.Remove(filePath); err != nil && verbose {
					fmt.Printf("⚠️  Warning: Cannot remove file %s: %v\n", filePath, err)
				}
			}
		}
	}

	// Final check: is the main directory gone?
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Success!
	}
	
	return fmt.Errorf("directory still exists after manual removal")
}

// forceRemoveWithSystemCommand uses system commands as last resort (Unix-like systems)
func forceRemoveWithSystemCommand(path string) error {
	// Only try this on Unix-like systems
	if filepath.Separator != '/' {
		return fmt.Errorf("system command removal not supported on this platform")
	}

	// Use rm -rf as last resort for cache directories
	cmd := exec.Command("rm", "-rf", path)
	return cmd.Run()
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func printStats(stats *CleanupStats) {
	fmt.Printf("📊 Cleanup Statistics:\n")
	fmt.Printf("   Projects processed: %d\n", stats.TotalProjects)
	fmt.Printf("   Cache items removed: %d\n", stats.TotalCacheItems)
	fmt.Printf("   Total space reclaimed: %s\n", formatBytes(stats.TotalSizeRemoved))
	fmt.Printf("   Processing time: %v\n", stats.ProcessingTime)
	if stats.ProcessingTime.Seconds() > 0 {
		fmt.Printf("   Average speed: %.2f MB/s\n",
			float64(stats.TotalSizeRemoved)/(1024*1024)/stats.ProcessingTime.Seconds())
	}
}

func listProjectTypes(config *Config) {
	fmt.Printf("📋 Supported Project Types (%d total):\n\n", len(config.ProjectTypes))

	for _, pt := range config.ProjectTypes {
		fmt.Printf("🔹 %s\n", pt.Name)
		fmt.Printf("   Indicators: %s\n", strings.Join(pt.Indicators, ", "))
		fmt.Printf("   Cache Directories: %s\n", strings.Join(pt.CacheConfig.Directories, ", "))

		if len(pt.CacheConfig.Files) > 0 {
			fmt.Printf("   Cache Files: %s\n", strings.Join(pt.CacheConfig.Files, ", "))
		}

		if len(pt.CacheConfig.Extensions) > 0 {
			fmt.Printf("   Cache Extensions: %s\n", strings.Join(pt.CacheConfig.Extensions, ", "))
		}

		fmt.Println()
	}

	fmt.Println("💡 Tip: Use --save-config to create a customizable configuration file")
}
