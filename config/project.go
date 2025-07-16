package config

import (
	"path/filepath"

	"github.com/core-stack/zetten-cli/internal/util"
)

type ProjectConfig struct {
	Name         string     `yaml:"name"`
	Version      string     `yaml:"version"`
	Dependencies Dependency `yaml:"dependencies"`
	Path         string     `yaml:"path"`

	DefaultProvider string `yaml:"default_provider"`

	dir string
}

func (p *ProjectConfig) Save() error {
	configPath := filepath.Join(p.dir, "zetten.yml")
	return util.SaveYAMLIndented(configPath, p)
}

func (p *ProjectConfig) AddDependency(name, version string) {
	p.Dependencies[name] = version
}
func (p *ProjectConfig) RemoveDependency(name string) {
	delete(p.Dependencies, name)
}

func (p *ProjectConfig) SetVersion(version string) {
	p.Version = version
}
