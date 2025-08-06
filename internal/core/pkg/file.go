package pkg

import (
	"path/filepath"
	"strings"

	"github.com/core-stack/zetten-cli/internal/util"
)

type PackageFile struct {
	Tag        string `yaml:"tag,omitempty"`
	Repository string `yaml:"repository"`
	Path       string `yaml:"-"`
}

func (f *PackageFile) Save() error {
	configPath := filepath.Join(f.Path)
	return util.SaveYAMLIndented(configPath, f)
}

func (c *PackageFile) GetName() string {
	if c.Repository == "" {
		return ""
	}
	parts := strings.Split(c.Repository, "/")
	lastPart := parts[len(parts)-1]
	return strings.TrimSuffix(lastPart, ".git")
}

func (c *PackageFile) GetTag() string {
	return c.Tag
}
func (c *PackageFile) GetRepository() string {
	return c.Repository
}
func (c *PackageFile) SetTag(tag string, autoSave bool) error {
	c.Tag = tag
	if autoSave {
		return c.Save()
	}
	return nil
}
func (c *PackageFile) SetRepository(repoURL string, autoSave bool) error {
	c.Repository = repoURL
	if autoSave {
		return c.Save()
	}
	return nil
}
func (c *PackageFile) UpdateFrom(other *PackageFile, autoSave bool) error {
	if other.Repository != "" {
		c.Repository = other.Repository
	}
	if other.Tag != "" {
		c.SetTag(other.Tag, false)
	}
	if autoSave {
		return c.Save()
	}
	return nil
}
