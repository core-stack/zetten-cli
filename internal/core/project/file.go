package project

import (
	"path/filepath"

	"github.com/core-stack/zetten-cli/internal/util"
)

type Dependency map[string]string

type ProjectFile struct {
	Name         string     `yaml:"name"`
	Version      string     `yaml:"version"`
	Dependencies Dependency `yaml:"dependencies"`
	PackagesPath string     `yaml:"packagesPath"`

	Path string `yaml:"-"`
}

func (f *ProjectFile) Save() error {
	configPath := filepath.Join(f.Path)
	return util.SaveYAMLIndented(configPath, f)
}

func (p *ProjectFile) AddDependency(url, version string, autoSave bool) error {
	if p.Dependencies == nil {
		p.Dependencies = make(map[string]string)
	}
	p.Dependencies[url] = version
	if autoSave {
		return p.Save()
	}
	return nil
}
func (p *ProjectFile) RemoveDependency(url string, autoSave bool) error {
	if p.Dependencies == nil {
		p.Dependencies = make(map[string]string)
	}
	delete(p.Dependencies, url)
	if autoSave {
		return p.Save()
	}
	return nil
}
func (p *ProjectFile) SetVersion(version string, autoSave bool) error {
	p.Version = version
	if autoSave {
		return p.Save()
	}
	return nil
}
