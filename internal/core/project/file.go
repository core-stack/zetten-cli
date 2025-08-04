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

func (p *ProjectFile) AddDependency(name, version string, autoSave bool) error {
	p.Dependencies[name] = version
	if autoSave {
		return p.Save()
	}
	return nil
}
func (p *ProjectFile) RemoveDependency(name string, autoSave bool) error {
	delete(p.Dependencies, name)
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
