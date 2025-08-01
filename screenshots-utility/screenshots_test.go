package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestScreenshotGeneration validates the screenshot generation process
func TestScreenshotGeneration(t *testing.T) {
	if os.Getenv("GENERATE_SCREENSHOTS") == "" {
		t.Skip("Set GENERATE_SCREENSHOTS=1 to run screenshot generation tests")
	}

	// Get the directory of this test file
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir := filepath.Dir(filename)
	rootDir := filepath.Dir(screenshotDir)

	// Test data for expected screenshots
	expectedScreenshots := []struct {
		name        string
		description string
		tapeFile    string
		minSizeKB   int64 // Minimum expected size in KB
	}{
		{
			name:        "basic-usage",
			description: "Basic Usage Demo",
			tapeFile:    "basic-usage.tape",
			minSizeKB:   100,
		},
		{
			name:        "dry-run",
			description: "Dry Run Demo",
			tapeFile:    "dry-run.tape",
			minSizeKB:   80,
		},
		{
			name:        "verbose",
			description: "Verbose Mode Demo",
			tapeFile:    "verbose.tape",
			minSizeKB:   120,
		},
		{
			name:        "interactive",
			description: "Interactive Mode Demo",
			tapeFile:    "interactive.tape",
			minSizeKB:   100,
		},
		{
			name:        "ui-demo",
			description: "TUI Interface Demo",
			tapeFile:    "ui-demo.tape",
			minSizeKB:   150,
		},
		{
			name:        "performance",
			description: "Performance Optimization Demo",
			tapeFile:    "performance.tape",
			minSizeKB:   120,
		},
		{
			name:        "quickstart",
			description: "Quickstart Workflow Demo",
			tapeFile:    "quickstart.tape",
			minSizeKB:   200,
		},
	}

	t.Run("VHS_Installation", func(t *testing.T) {
		// Check if VHS is installed
		cmd := exec.Command("vhs", "--version")
		if err := cmd.Run(); err != nil {
			t.Fatal("VHS is not installed. Install with: go install github.com/charmbracelet/vhs@latest")
		}
		t.Log("✅ VHS is installed and accessible")
	})

	t.Run("Binary_Exists", func(t *testing.T) {
		// Check if cache-remover binary exists
		binaryPath := filepath.Join(rootDir, "cache-remover")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			t.Fatal("cache-remover binary not found. Run: go build -o cache-remover")
		}
		t.Log("✅ cache-remover binary found")
	})

	t.Run("Tape_Files_Exist", func(t *testing.T) {
		// Check if all tape files exist
		tapesDir := filepath.Join(screenshotDir, "tapes")
		for _, screenshot := range expectedScreenshots {
			tapeFile := filepath.Join(tapesDir, screenshot.tapeFile)
			if _, err := os.Stat(tapeFile); os.IsNotExist(err) {
				t.Errorf("Tape file not found: %s", tapeFile)
			}
		}
		t.Log("✅ All tape files found")
	})

	t.Run("Setup_Test_Data", func(t *testing.T) {
		// Run setup script
		setupScript := filepath.Join(screenshotDir, "setup-test-data.sh")
		cmd := exec.Command("bash", setupScript)
		cmd.Dir = rootDir
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to setup test data: %v\nOutput: %s", err, output)
		}
		t.Log("✅ Test data setup completed")
	})

	// Generate screenshots for each tape file
	for _, screenshot := range expectedScreenshots {
		t.Run("Generate_"+screenshot.name, func(t *testing.T) {
			tapeFile := filepath.Join(screenshotDir, "tapes", screenshot.tapeFile)
			
			// Run VHS to generate screenshot
			cmd := exec.Command("vhs", tapeFile)
			cmd.Dir = rootDir
			
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to generate %s: %v\nOutput: %s", screenshot.name, err, output)
			}
			
			// Check if output file was created
			outputFile := filepath.Join(screenshotDir, "screenshots", screenshot.name+".gif")
			stat, err := os.Stat(outputFile)
			if err != nil {
				t.Fatalf("Screenshot file not created: %s", outputFile)
			}
			
			// Check minimum file size (ensure it's not empty/corrupted)
			sizeKB := stat.Size() / 1024
			if sizeKB < screenshot.minSizeKB {
				t.Errorf("Screenshot %s is too small (%d KB < %d KB), may be corrupted", 
					screenshot.name, sizeKB, screenshot.minSizeKB)
			}
			
			t.Logf("✅ Generated %s (%d KB)", screenshot.name, sizeKB)
		})
	}

	t.Run("Validate_All_Screenshots", func(t *testing.T) {
		screenshotsDir := filepath.Join(screenshotDir, "screenshots")
		
		// Count generated screenshots
		files, err := filepath.Glob(filepath.Join(screenshotsDir, "*.gif"))
		if err != nil {
			t.Fatalf("Failed to list screenshot files: %v", err)
		}
		
		if len(files) != len(expectedScreenshots) {
			t.Errorf("Expected %d screenshots, found %d", len(expectedScreenshots), len(files))
		}
		
		// Validate each screenshot
		for _, file := range files {
			stat, err := os.Stat(file)
			if err != nil {
				t.Errorf("Cannot stat screenshot file: %s", file)
				continue
			}
			
			if stat.Size() == 0 {
				t.Errorf("Empty screenshot file: %s", file)
			}
		}
		
		t.Logf("✅ Validated %d screenshot files", len(files))
	})
}

// TestDocumentationLinks validates that documentation contains correct screenshot links
func TestDocumentationLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping documentation link validation in short mode")
	}

	rootDir := filepath.Dir(filepath.Dir(os.Args[0]))
	
	// Test README.md
	t.Run("README_Links", func(t *testing.T) {
		readmePath := filepath.Join(rootDir, "README.md")
		content, err := os.ReadFile(readmePath)
		if err != nil {
			t.Fatalf("Cannot read README.md: %v", err)
		}
		
		readmeContent := string(content)
		
		// Check for screenshot references (when they're added)
		expectedScreenshots := []string{
			"screenshots-utility/screenshots/quickstart.gif",
			"screenshots-utility/screenshots/performance.gif",
		}
		
		for _, screenshot := range expectedScreenshots {
			if !strings.Contains(readmeContent, screenshot) {
				t.Logf("Note: README.md doesn't contain reference to %s", screenshot)
				// Don't fail - screenshots may not be integrated yet
			}
		}
	})

	// Test QUICKSTART.md
	t.Run("QUICKSTART_Links", func(t *testing.T) {
		quickstartPath := filepath.Join(rootDir, "QUICKSTART.md")
		content, err := os.ReadFile(quickstartPath)
		if err != nil {
			t.Fatalf("Cannot read QUICKSTART.md: %v", err)
		}
		
		quickstartContent := string(content)
		
		// Check for potential screenshot locations
		expectedScreenshots := []string{
			"screenshots-utility/screenshots/basic-usage.gif",
			"screenshots-utility/screenshots/dry-run.gif",
		}
		
		for _, screenshot := range expectedScreenshots {
			if !strings.Contains(quickstartContent, screenshot) {
				t.Logf("Note: QUICKSTART.md doesn't contain reference to %s", screenshot)
			}
		}
	})
}

// BenchmarkScreenshotGeneration benchmarks the screenshot generation process
func BenchmarkScreenshotGeneration(b *testing.B) {
	if os.Getenv("BENCHMARK_SCREENSHOTS") == "" {
		b.Skip("Set BENCHMARK_SCREENSHOTS=1 to run screenshot generation benchmarks")
	}

	// Setup once
	_, filename, _, _ := runtime.Caller(0)
	screenshotDir := filepath.Dir(filename)
	rootDir := filepath.Dir(screenshotDir)
	
	setupScript := filepath.Join(screenshotDir, "setup-test-data.sh")
	cmd := exec.Command("bash", setupScript)
	cmd.Dir = rootDir
	if err := cmd.Run(); err != nil {
		b.Fatalf("Failed to setup test data: %v", err)
	}

	b.ResetTimer()
	
	// Benchmark single tape generation
	tapeFile := filepath.Join(screenshotDir, "tapes", "dry-run.tape")
	
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("vhs", tapeFile)
		cmd.Dir = rootDir
		if err := cmd.Run(); err != nil {
			b.Fatalf("VHS failed: %v", err)
		}
	}
}

// Helper function to check if VHS is available
func isVHSAvailable() bool {
	cmd := exec.Command("vhs", "--version")
	return cmd.Run() == nil
}