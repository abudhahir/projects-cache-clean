package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	Name        string
	Indicators  []string
	CacheConfig CacheConfig
}

var projectTypes = []ProjectType{
	{
		Name:       "Node.js",
		Indicators: []string{"package.json", "yarn.lock", "package-lock.json"},
		CacheConfig: CacheConfig{
			Directories: []string{"node_modules", "dist", "build", ".next", ".nuxt", "coverage"},
			Files:       []string{},
			Extensions:  []string{},
		},
	},
	{
		Name:       "Python",
		Indicators: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
		CacheConfig: CacheConfig{
			Directories: []string{"__pycache__", ".pytest_cache", "dist", "build", ".mypy_cache", ".tox", "venv", ".venv"},
			Files:       []string{},
			Extensions:  []string{".pyc", ".pyo"},
		},
	},
	{
		Name:       "Java/Maven",
		Indicators: []string{"pom.xml"},
		CacheConfig: CacheConfig{
			Directories: []string{"target"},
			Files:       []string{},
			Extensions:  []string{},
		},
	},
	{
		Name:       "Gradle",
		Indicators: []string{"build.gradle", "build.gradle.kts"},
		CacheConfig: CacheConfig{
			Directories: []string{"build", ".gradle"},
			Files:       []string{},
			Extensions:  []string{},
		},
	},
	{
		Name:       "Go",
		Indicators: []string{"go.mod", "go.sum"},
		CacheConfig: CacheConfig{
			Directories: []string{"vendor"},
			Files:       []string{},
			Extensions:  []string{},
		},
	},
	{
		Name:       "Rust",
		Indicators: []string{"Cargo.toml"},
		CacheConfig: CacheConfig{
			Directories: []string{"target"},
			Files:       []string{},
			Extensions:  []string{},
		},
	},
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
	var (
		rootDir    = flag.String("dir", ".", "Root directory to scan for projects")
		dryRun     = flag.Bool("dry-run", false, "Show what would be removed without actually removing")
		workers    = flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
		verbose    = flag.Bool("verbose", false, "Verbose output")
		maxDepth   = flag.Int("max-depth", 10, "Maximum directory depth to scan")
		interactive = flag.Bool("interactive", false, "Ask for confirmation before removing each cache")
		ui         = flag.Bool("ui", false, "Launch interactive TUI mode")
		webui      = flag.Bool("web", false, "Launch web UI mode")
		port       = flag.Int("port", 8080, "Port for web UI")
	)
	flag.Parse()

	// Launch web UI if requested
	if *webui {
		fmt.Println("🌐 Launching Web UI Cache Remover...")
		if err := runWebUI(*rootDir, *port); err != nil {
			fmt.Printf("Error running web UI: %v\n", err)
			os.Exit(1)
		}
		return
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

func findProjects(rootDir string, maxDepth int, verbose bool) []string {
	var projects []string
	var mu sync.Mutex

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		depth := strings.Count(strings.TrimPrefix(path, rootDir), string(os.PathSeparator))
		if depth > maxDepth {
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
	for _, projectType := range projectTypes {
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
	for _, projectType := range projectTypes {
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

	for _, dir := range config.Directories {
		dirPath := filepath.Join(projectPath, dir)
		if size := getDirSize(dirPath); size > 0 {
			items = append(items, CacheItem{
				Path: dirPath,
				Size: size,
				Type: "directory",
			})
		}
	}

	for _, file := range config.Files {
		filePath := filepath.Join(projectPath, file)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			items = append(items, CacheItem{
				Path: filePath,
				Size: info.Size(),
				Type: "file",
			})
		}
	}

	if len(config.Extensions) > 0 {
		filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

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
			return nil
		})
	}

	return items
}

func getDirSize(dirPath string) int64 {
	var size int64
	filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
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
		if err := os.RemoveAll(item.Path); err != nil {
			if verbose {
				fmt.Printf("❌ Failed to remove %s: %v\n", item.Path, err)
			}
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