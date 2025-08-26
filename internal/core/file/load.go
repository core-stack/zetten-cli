package file

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func Load[T any](path string) (*T, error) {
	var cfg T

	if path == "" {
		path = "."
	}
	if _, err := os.Stat(path); err != nil {
		return &cfg, fmt.Errorf("file not found")
	}

	data, err := os.ReadFile(path)
	fmt.Println(path)
	if err != nil {
		return &cfg, fmt.Errorf("❌ Failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &cfg, fmt.Errorf("❌ Failed to parse config: %w", err)
	}
	return &cfg, nil
}
