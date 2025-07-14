package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Dependency map[string]string

type ProjectConfig struct {
	Name         string     `yaml:"name"`
	Version      string     `yaml:"version"`
	Dependencies Dependency `yaml:"dependencies"`
	Path         string     `yaml:"path"`
}

type PackageConfig struct {
	Name         string     `yaml:"name"`
	Version      string     `yaml:"version"`
	Private      bool       `yaml:"private"`
	Provider     string     `yaml:"provider"` // github, gitlab, bitbucket
	Dependencies Dependency `yaml:"dependencies"`
	Repository   string     `yaml:"repository"`
}

func LoadProjectConfig(dir string) (*ProjectConfig, error) {
	if dir == "" {
		dir = "."
	}
	configPath := filepath.Join(dir, "zetten.yml")
	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("⚠️ Configuration file not found. Navigate to the root directory and run `zetten init`")
	}
	var cfg ProjectConfig

	data, err := os.ReadFile(configPath)
	if err != nil {
		return &cfg, fmt.Errorf("❌ Failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &cfg, fmt.Errorf("❌ Failed to parse config: %w", err)
	}

	return &cfg, nil
}
