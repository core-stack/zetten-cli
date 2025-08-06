package root

import (
	"path/filepath"

	"slices"

	"github.com/core-stack/zetten-cli/internal/util"
)

type RootFile struct {
	ZettenProjects []string   `yaml:"zettenProjects"`
	Path           string     `yaml:"-"`
	Mirror         [][]string `yaml:"mirror"`
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

func (r *RootFile) AddMirror(paths []string, autoSave bool) error {
	r.Mirror = append(r.Mirror, paths)
	if autoSave {
		return r.Save()
	}
	return nil
}

func (r *RootFile) RemoveMirror(path string, autoSave bool) error {
	for i, p := range r.Mirror {
		for j, m := range p {
			if m == path {
				r.Mirror[i] = slices.Delete(r.Mirror[i], j, j+1)
			}
		}
	}
	if autoSave {
		return r.Save()
	}
	return nil
}
