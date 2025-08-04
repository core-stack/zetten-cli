package root

import (
	"path/filepath"

	"github.com/core-stack/zetten-cli/internal/util"
)

type RootFile struct {
	ZettenProjects []string `yaml:"zettenProjects"`
	Path           string   `yaml:"-"`
}

func (f *RootFile) Save() error {
	configPath := filepath.Join(f.Path)
	return util.SaveYAMLIndented(configPath, f)
}

func (r *RootFile) AddProject(dir string, autoSave bool) error {
	r.ZettenProjects = append(r.ZettenProjects, dir)
	if autoSave {
		return r.Save()
	}
	return nil
}

func (r *RootFile) RemoveProject(dir string, autoSave bool) error {
	for i, p := range r.ZettenProjects {
		if p == dir {
			r.ZettenProjects = append(r.ZettenProjects[:i], r.ZettenProjects[i+1:]...)
			if autoSave {
				return r.Save()
			}
			return nil
		}
	}
	return nil
}
