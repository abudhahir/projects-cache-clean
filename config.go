package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ProjectTypes []ProjectType `json:"project_types"`
	Settings     Settings      `json:"settings"`
}

type Settings struct {
	MaxDepth       int    `json:"max_depth"`
	DefaultWorkers int    `json:"default_workers"`
	LogLevel       string `json:"log_level"`
}

var (
	globalConfig *Config
	configPaths  = []string{
		"config.json",               // Current directory
		"cache-remover-config.json", // Current directory with app prefix
		filepath.Join(os.Getenv("HOME"), ".cache-remover", "config.json"), // User home
		"/etc/cache-remover/config.json",                                  // System-wide
	}
)

// loadConfig loads configuration from various possible locations
func loadConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// Try to load from each possible config path
	for _, configPath := range configPaths {
		if config, err := loadConfigFromFile(configPath); err == nil {
			globalConfig = config
			fmt.Printf("📄 Loaded configuration from: %s\n", configPath)
			return globalConfig, nil
		}
	}

	// Fallback to default configuration
	config := getDefaultConfig()
	globalConfig = &config
	fmt.Println("📄 Using default configuration")
	return globalConfig, nil
}

func loadConfigFromFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid JSON in config file %s: %v", configPath, err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %v", configPath, err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if len(config.ProjectTypes) == 0 {
		return fmt.Errorf("no project types defined")
	}

	for i, pt := range config.ProjectTypes {
		if pt.Name == "" {
			return fmt.Errorf("project type %d has empty name", i)
		}
		if len(pt.Indicators) == 0 {
			return fmt.Errorf("project type '%s' has no indicators", pt.Name)
		}
	}

	if config.Settings.MaxDepth <= 0 {
		config.Settings.MaxDepth = 10
	}
	if config.Settings.DefaultWorkers <= 0 {
		config.Settings.DefaultWorkers = 4
	}
	if config.Settings.LogLevel == "" {
		config.Settings.LogLevel = "info"
	}

	return nil
}

func getDefaultConfig() Config {
	return Config{
		ProjectTypes: []ProjectType{
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
			{
				Name:       "Angular",
				Indicators: []string{"angular.json"},
				CacheConfig: CacheConfig{
					Directories: []string{"node_modules", "dist", ".angular"},
					Files:       []string{},
					Extensions:  []string{},
				},
			},
			{
				Name:       "Flutter",
				Indicators: []string{"pubspec.yaml"},
				CacheConfig: CacheConfig{
					Directories: []string{"build", ".dart_tool"},
					Files:       []string{},
					Extensions:  []string{},
				},
			},
			{
				Name:       "Swift/iOS",
				Indicators: []string{"Package.swift", "*.xcodeproj", "*.xcworkspace"},
				CacheConfig: CacheConfig{
					Directories: []string{"build", "DerivedData", ".build"},
					Files:       []string{},
					Extensions:  []string{},
				},
			},
		},
		Settings: Settings{
			MaxDepth:       10,
			DefaultWorkers: 4,
			LogLevel:       "info",
		},
	}
}

// saveDefaultConfig creates a default config file in the current directory
func saveDefaultConfig() error {
	config := getDefaultConfig()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("cache-remover-config.json", data, 0644)
}
