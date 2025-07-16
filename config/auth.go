package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var AUTH_FILE_NAME = "zetten-auth.yml"

type AuthConfig struct {
	Method      string `yaml:"method"`
	Credentials string `yaml:"credentials"`
}

type AuthMap map[string]AuthConfig

func LoadAuthForHost(host string) (*AuthConfig, error) {
	paths := []string{
		AUTH_FILE_NAME, // local
		filepath.Join(os.Getenv("HOME"), fmt.Sprintf(".zetten/%s", AUTH_FILE_NAME)), // global
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var auths AuthMap
		if err := yaml.Unmarshal(data, &auths); err != nil {
			return nil, fmt.Errorf("invalid YAML in %s: %w", path, err)
		}

		if config, ok := auths[host]; ok {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("no auth config found for host: %s", host)
}
