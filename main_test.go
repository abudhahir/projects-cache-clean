package main

import (
	"os"
	"path/filepath"
	"testing"
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