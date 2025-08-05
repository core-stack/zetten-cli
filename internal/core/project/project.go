package project

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/core-stack/zetten-cli/internal/core/file"
	"github.com/core-stack/zetten-cli/internal/core/root"
	"github.com/core-stack/zetten-cli/internal/util"
)

type ProjectConfig struct {
	ProjectFile `yaml:",inline"`

	Root root.IRootConfig
}

func (p *ProjectConfig) CopyFromRoot(url string) error {
	path := strings.TrimSuffix(util.ExtractPathFromURL(url), ".git")
	destination := filepath.Join(p.PackagesPath, path)
	return p.Root.CopyRootFiles(url, destination, []string{".git"})
}

func (p *ProjectConfig) Install(url, tagOrBranch string) error {
	if url == "" {
		return errors.New("url is required")
	}
	if tagOrBranch == "" {
		return errors.New("tag or branch is required")
	}

	_, err := p.Root.OpenOrClonePackage(url)
	if err != nil {
		return err
	}
	err = p.Root.Checkout(url, tagOrBranch)
	if err != nil {
		return err
	}

	err = p.CopyFromRoot(url)
	if err != nil {
		return err
	}

	if err = p.AddDependency(url, tagOrBranch, true); err != nil {
		return errors.New("error saving new dependency")
	}

	return nil
}

func (p *ProjectConfig) Uninstall(urls []string) error {
	if p.Dependencies == nil || len(p.Dependencies) == 0 {
		return errors.New("no dependencies found")
	}
	if len(urls) == 0 {
		return errors.New("no urls provided")
	}
	for _, url := range urls {
		if url == "" {
			continue
		}
		err := os.RemoveAll(filepath.Join(p.PackagesPath, util.ExtractPathFromURL(url)))
		if err != nil {
			return err
		}
		p.RemoveDependency(url, false)
	}
	err := p.cleanPackageFolders()
	if err != nil {
		return err
	}
	return p.Save()
}

func (p *ProjectConfig) cleanPackageFolders() error {
	var emptyFolders []string
	files, err := os.ReadDir(p.PackagesPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			entries, err := os.ReadDir(filepath.Join(p.PackagesPath, f.Name()))
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				emptyFolders = append(emptyFolders, filepath.Join(p.PackagesPath, f.Name()))
			}
		}
	}

	for _, folder := range emptyFolders {
		err := os.RemoveAll(folder)
		if err != nil {
			continue
		}
	}
	return nil
}

func (p *ProjectConfig) Sync() error {
	for url, version := range p.Dependencies {
		if err := p.Install(url, version); err != nil {
			return err
		}
	}
	return nil
}

func LoadProjectConfig(path string) (*ProjectConfig, error) {
	cfg, err := file.Load[ProjectConfig](path)
	if err != nil {
		return nil, err
	}
	if root, err := root.LoadRootConfig(); err != nil {
		return nil, err
	} else {
		cfg.Root = root
	}
	cfg.Path = path
	return cfg, nil
}

func NewProjectConfig(path, name, version, packagesPath string) (*ProjectConfig, error) {
	root, err := root.LoadRootConfig()
	if err != nil {
		return nil, err
	}
	cfg := ProjectConfig{
		ProjectFile: ProjectFile{
			Name:         name,
			Version:      version,
			PackagesPath: packagesPath,
			Path:         path,
			Dependencies: Dependency{},
		},
		Root: root,
	}
	err = cfg.Save()
	return &cfg, err
}
