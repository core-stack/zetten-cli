package config

type PackageConfig struct {
	Name         string     `yaml:"name"`
	Version      string     `yaml:"version"`
	Private      bool       `yaml:"private"`
	Provider     string     `yaml:"provider"` // github, gitlab, bitbucket
	Dependencies Dependency `yaml:"dependencies"`
	Repository   string     `yaml:"repository"`
}
