package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/stretchr/testify/assert"
)

func TestNewAndLoadProjectConfig(t *testing.T) {
	tmp := t.TempDir()
	configPath := filepath.Join(tmp, "project.yaml")

	packagesPath := filepath.Join(tmp, "pkgs")
	err := os.MkdirAll(packagesPath, 0755)
	assert.NoError(t, err)

	// Criar nova config
	cfg, err := project.NewProjectConfig(configPath, "my-app", "0.1.0", packagesPath)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "my-app", cfg.Name)
	assert.Equal(t, "0.1.0", cfg.Version)

	// Carregar config salva
	loaded, err := project.LoadProjectConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, "my-app", loaded.Name)
	assert.Equal(t, packagesPath, loaded.PackagesPath)
	assert.Equal(t, configPath, loaded.Path)
}

func TestInstall_Success(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")
	pkgs := filepath.Join(tmp, "packages")

	cfg := &project.ProjectConfig{
		ProjectFile: project.ProjectFile{
			Name:         "test",
			Version:      "1.0.0",
			Path:         path,
			PackagesPath: pkgs,
			Dependencies: project.Dependency{},
		},
		Root: &MockRootConfig{},
	}

	err := cfg.Install("github.com/user/repo", "v1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", cfg.Dependencies["github.com/user/repo"])
}

func TestInstall_MissingURL(t *testing.T) {
	cfg := &project.ProjectConfig{}
	err := cfg.Install("", "v1.0.0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "url is required")
}

func TestInstall_CheckoutFails(t *testing.T) {
	cfg := &project.ProjectConfig{
		ProjectFile: project.ProjectFile{
			PackagesPath: t.TempDir(),
			Dependencies: project.Dependency{},
		},
		Root: &MockRootConfig{},
	}
	err := cfg.Install("github.com/user/repo", "error")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checkout failed")
}
func TestRemove_Success(t *testing.T) {
	tmp := t.TempDir()
	pkgPath := filepath.Join(tmp, "packages", "github.com/user/repo")
	os.MkdirAll(pkgPath, 0755)

	cfg := &project.ProjectConfig{
		ProjectFile: project.ProjectFile{
			PackagesPath: filepath.Join(tmp, "packages"),
			Dependencies: project.Dependency{"github.com/user/repo": "v1.0.0"},
			Path:         filepath.Join(tmp, "project.yaml"),
		},
	}

	err := cfg.Uninstall([]string{"github.com/user/repo"})
	assert.NoError(t, err)
	assert.NotContains(t, cfg.Dependencies, "github.com/user/repo")
	assert.NoDirExists(t, pkgPath)
}

func TestSync_Success(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")
	pkgs := filepath.Join(tmp, "packages")

	cfg := &project.ProjectConfig{
		ProjectFile: project.ProjectFile{
			Path:         path,
			PackagesPath: pkgs,
			Dependencies: project.Dependency{
				"github.com/user/repo1": "v1.0.0",
				"github.com/user/repo2": "v2.0.0",
			},
		},
		Root: &MockRootConfig{},
	}

	err := cfg.Sync()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cfg.Dependencies))
}
