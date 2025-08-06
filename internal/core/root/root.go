package root

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/core-stack/zetten-cli/internal/core/file"
	"github.com/core-stack/zetten-cli/internal/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var home, err = os.UserHomeDir()
var DEFAULT_ROOT_PATH = filepath.Join(home, ".zetten")
var DEFAULT_ROOT_CONFIG_PATH = filepath.Join(DEFAULT_ROOT_PATH, "config.yml")
var DEFAULT_ROOT_PACKAGES_PATH = filepath.Join(DEFAULT_ROOT_PATH, "packages")

type IRootConfig interface {
	BuildRootPackagePath(url string) string
	HasPackage(url string) bool
	OpenOrClonePackage(url string) (*git.Repository, error)
	Checkout(url, tag string) (*git.Repository, error)
	CopyRootFiles(url, destination string, ignore []string) error
	Promote(url, tag, newTag, packageDir string) error
}
type RootConfig struct {
	RootFile `yaml:",inline"`
}

func (r *RootConfig) BuildRootPackagePath(url string) string {
	return filepath.Join(DEFAULT_ROOT_PACKAGES_PATH, strings.TrimSuffix(util.ExtractPathFromURL(url), ".git"))
}

func (r *RootConfig) HasPackage(url string) bool {
	_, err := os.Stat(r.BuildRootPackagePath(url))
	return err == nil
}

func (r *RootConfig) OpenOrClonePackage(url string) (*git.Repository, error) {
	destination := r.BuildRootPackagePath(url)
	if r.HasPackage(url) {
		return git.PlainOpen(destination)
	}

	repo, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return repo, err
	}
	return repo, nil
}

func (r *RootConfig) Checkout(url, tag string) (*git.Repository, error) {
	if tag == "" {
		return nil, fmt.Errorf("No tag specified")
	}
	if url == "" {
		return nil, fmt.Errorf("No url specified")
	}
	repo, err := r.OpenOrClonePackage(url)
	if err != nil {
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	err = wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(tag)})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RootConfig) Promote(url, baseTag, newTag, packageDir string) error {
	repo, err := r.Checkout(url, baseTag)
	if err != nil {
		return err
	}

	srcDir := r.BuildRootPackagePath(url)
	err = util.CopyDir(packageDir, srcDir, []string{})
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}

	commitHash, err := wt.Commit(fmt.Sprintf("Promote to tag %s", newTag), nil)
	if err != nil {
		return err
	}

	_, err = repo.CreateTag(newTag, commitHash, nil)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Changes promoted and tagged as %s\n", newTag)
	return nil
}

func (r *RootConfig) CopyRootFiles(url, packagesDir string, ignore []string) error {
	ignore = append(ignore, ".git")
	srcDir := r.BuildRootPackagePath(url)
	return util.CopyDir(srcDir, packagesDir, ignore)
}

func LoadRootConfig() (*RootConfig, error) {
	if _, err := os.Stat(DEFAULT_ROOT_PATH); os.IsNotExist(err) {
		if err := os.MkdirAll(DEFAULT_ROOT_PATH, os.ModePerm); err != nil {
			return nil, fmt.Errorf("❌ Failed to create root config directory: %w", err)
		}
	}

	if _, err := os.Stat(DEFAULT_ROOT_CONFIG_PATH); os.IsNotExist(err) {
		defaultRoot := &RootFile{
			Path:           DEFAULT_ROOT_CONFIG_PATH,
			ZettenProjects: []string{},
		}
		if err := defaultRoot.Save(); err != nil {
			return nil, fmt.Errorf("❌ Failed to create default root config: %w", err)
		}
	}
	cfg, err := file.Load[RootConfig](DEFAULT_ROOT_CONFIG_PATH)
	if err != nil {
		return nil, err
	}
	cfg.Path = DEFAULT_ROOT_CONFIG_PATH
	return cfg, nil
}
