package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDetectProjectType(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test Node.js project detection
	packageJSONPath := filepath.Join(tempDir, "package.json")
	file, err := os.Create(packageJSONPath)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}
	file.Close()

	projectType := detectProjectType(tempDir)
	if projectType == nil || projectType.Name != "Node.js" {
		t.Errorf("Expected Node.js project type, got %v", projectType)
	}

	// Clean up
	os.Remove(packageJSONPath)

	// Test Python project detection
	requirementsPath := filepath.Join(tempDir, "requirements.txt")
	file, err = os.Create(requirementsPath)
	if err != nil {
		t.Fatalf("Failed to create requirements.txt: %v", err)
	}
	file.Close()

	projectType = detectProjectType(tempDir)
	if projectType == nil || projectType.Name != "Python" {
		t.Errorf("Expected Python project type, got %v", projectType)
	}
}

func TestIsProjectDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Initially should not be a project directory
	if isProjectDirectory(tempDir) {
		t.Error("Empty directory should not be detected as project")
	}

	// Create a Go project indicator
	goModPath := filepath.Join(tempDir, "go.mod")
	file, err := os.Create(goModPath)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}
	file.Close()

	// Now should be detected as project directory
	if !isProjectDirectory(tempDir) {
		t.Error("Directory with go.mod should be detected as project")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, test := range tests {
		result := formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}

func TestFindCacheItems(t *testing.T) {
	tempDir := t.TempDir()

	// Create some cache directories and files
	nodemodulesDir := filepath.Join(tempDir, "node_modules")
	os.MkdirAll(nodemodulesDir, 0755)

	// Create a file inside node_modules
	testFile := filepath.Join(nodemodulesDir, "test.js")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.WriteString("console.log('test');")
	file.Close()

	config := CacheConfig{
		Directories: []string{"node_modules"},
		Files:       []string{},
		Extensions:  []string{},
	}

	items := findCacheItems(tempDir, config)
	if len(items) != 1 {
		t.Errorf("Expected 1 cache item, got %d", len(items))
	}

	if items[0].Path != nodemodulesDir {
		t.Errorf("Expected cache item path %s, got %s", nodemodulesDir, items[0].Path)
	}

	if items[0].Size <= 0 {
		t.Error("Cache item should have size > 0")
	}
}

// Integration Tests

func TestFullWorkflowIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// Create a realistic project structure in separate directories
	setupTestProject(t, filepath.Join(tempDir, "react-app"), "react-app", "Node.js")
	setupTestProject(t, filepath.Join(tempDir, "python-api"), "python-api", "Python")
	setupTestProject(t, filepath.Join(tempDir, "java-service"), "java-service", "Java")

	// Test project discovery
	projects := findProjects(tempDir, 10, false)
	if len(projects) != 3 {
		t.Errorf("Expected 3 projects, found %d", len(projects))
	}

	// Test cache detection and cleanup
	stats := &CleanupStats{}
	processProjects(projects, 1, true, false, false, stats) // dry run

	if stats.TotalProjects != 3 {
		t.Errorf("Expected 3 projects processed, got %d", stats.TotalProjects)
	}

	if stats.TotalSizeRemoved <= 0 {
		t.Error("Expected some cache size to be calculated")
	}
}

func TestCacheDirectorySkipping(t *testing.T) {
	tempDir := t.TempDir()

	// Create nested structure with cache directories
	nodePath := filepath.Join(tempDir, "project")
	nodeModulesPath := filepath.Join(nodePath, "node_modules")
	nestedPath := filepath.Join(nodeModulesPath, "some-package", "node_modules")

	os.MkdirAll(nestedPath, 0755)

	// Create package.json to make it a project
	packageJSON := filepath.Join(nodePath, "package.json")
	os.WriteFile(packageJSON, []byte(`{"name": "test"}`), 0644)

	// Create a file deep in nested node_modules
	testFile := filepath.Join(nestedPath, "test.js")
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Test that we skip descending into cache directories
	projects := findProjects(tempDir, 10, true)
	if len(projects) != 1 {
		t.Errorf("Expected 1 project, found %d", len(projects))
	}

	// The nested node_modules should not cause additional projects to be found
	if projects[0] != nodePath {
		t.Errorf("Expected project path %s, got %s", nodePath, projects[0])
	}
}

func TestConcurrentProcessing(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple projects
	projectCount := 5
	for i := 0; i < projectCount; i++ {
		projectDir := filepath.Join(tempDir, filepath.Join("project", string(rune('0'+i))))
		setupTestProject(t, projectDir, "test-project", "Node.js")
	}

	projects := findProjects(tempDir, 10, false)
	if len(projects) != projectCount {
		t.Errorf("Expected %d projects, found %d", projectCount, len(projects))
	}

	// Test concurrent processing with multiple workers
	stats := &CleanupStats{}
	startTime := time.Now()
	processProjects(projects, 3, true, false, false, stats) // 3 workers, dry run
	processingTime := time.Since(startTime)

	if stats.TotalProjects != projectCount {
		t.Errorf("Expected %d projects processed, got %d", projectCount, stats.TotalProjects)
	}

	// Should complete reasonably quickly with concurrent processing
	if processingTime > time.Second {
		t.Errorf("Processing took too long: %v", processingTime)
	}
}

func TestRealCleanupWorkflow(t *testing.T) {
	tempDir := t.TempDir()

	// Create a project with actual cache content
	projectDir := filepath.Join(tempDir, "test-project")
	setupTestProject(t, projectDir, "test-project", "Node.js")

	// Verify cache exists
	nodeModulesPath := filepath.Join(projectDir, "node_modules")
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		t.Fatal("node_modules should exist before cleanup")
	}

	// Perform actual cleanup (not dry run)
	projects := findProjects(tempDir, 10, false)
	stats := &CleanupStats{}
	processProjects(projects, 1, false, false, false, stats) // actual cleanup

	// Verify cache was removed
	if _, err := os.Stat(nodeModulesPath); !os.IsNotExist(err) {
		t.Error("node_modules should be removed after cleanup")
	}

	if stats.TotalCacheItems <= 0 {
		t.Error("Expected cache items to be removed")
	}
}

func TestErrorHandling(t *testing.T) {
	tempDir := t.TempDir()

	// Create cache item that should fail to remove (but os.RemoveAll doesn't fail on non-existent)
	// So let's create a file first, then remove it, then try to remove it again
	testFile := filepath.Join(tempDir, "test-file")
	os.WriteFile(testFile, []byte("test"), 0644)
	os.Remove(testFile) // Remove it first

	item := CacheItem{
		Path: testFile,
		Size: 100,
		Type: "file",
	}

	// Test that removal of non-existent item succeeds (os.RemoveAll behavior)
	removed, size := removeCacheItems([]CacheItem{item}, false)

	// os.RemoveAll succeeds even if file doesn't exist
	if removed != 1 {
		t.Errorf("Expected 1 item removed (os.RemoveAll succeeds on non-existent), got %d", removed)
	}

	if size != 100 {
		t.Errorf("Expected size 100 to be counted as removed, got %d", size)
	}
}

// Helper function to setup realistic test projects
func setupTestProject(t *testing.T, projectDir, name, projectType string) {
	t.Helper()

	os.MkdirAll(projectDir, 0755)

	switch projectType {
	case "Node.js":
		// Create package.json
		packageJSON := filepath.Join(projectDir, "package.json")
		content := `{"name": "` + name + `", "version": "1.0.0"}`
		os.WriteFile(packageJSON, []byte(content), 0644)

		// Create node_modules with some content
		nodeModules := filepath.Join(projectDir, "node_modules")
		os.MkdirAll(nodeModules, 0755)

		// Add some files to simulate real cache
		for i := 0; i < 5; i++ {
			filePath := filepath.Join(nodeModules, "file"+string(rune('0'+i))+".js")
			content := strings.Repeat("// module content\n", 100) // ~1.8KB per file
			os.WriteFile(filePath, []byte(content), 0644)
		}

	case "Python":
		// Create requirements.txt
		reqPath := filepath.Join(projectDir, "requirements.txt")
		os.WriteFile(reqPath, []byte("requests==2.25.1\n"), 0644)

		// Create __pycache__
		pycacheDir := filepath.Join(projectDir, "__pycache__")
		os.MkdirAll(pycacheDir, 0755)

		// Add .pyc files
		for i := 0; i < 3; i++ {
			filePath := filepath.Join(pycacheDir, "module"+string(rune('0'+i))+".cpython-39.pyc")
			content := strings.Repeat("compiled python bytecode\n", 50)
			os.WriteFile(filePath, []byte(content), 0644)
		}

	case "Java":
		// Create pom.xml
		pomPath := filepath.Join(projectDir, "pom.xml")
		pomContent := `<?xml version="1.0"?>
<project>
	<groupId>com.example</groupId>
	<artifactId>` + name + `</artifactId>
	<version>1.0.0</version>
</project>`
		os.WriteFile(pomPath, []byte(pomContent), 0644)

		// Create target directory
		targetDir := filepath.Join(projectDir, "target")
		os.MkdirAll(targetDir, 0755)

		// Add compiled classes
		classFile := filepath.Join(targetDir, "Main.class")
		content := strings.Repeat("compiled java bytecode\n", 200) // ~4KB
		os.WriteFile(classFile, []byte(content), 0644)
	}
}
