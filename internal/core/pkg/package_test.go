package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/pkg"
	"github.com/stretchr/testify/assert"
)

func TestLoadPackageConfig_NotFound(t *testing.T) {
	cfg, err := pkg.LoadPackageConfig("notfound.yaml")
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Configuration file not found")
}

func TestLoadPackageConfig_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "invalid.yaml")
	os.WriteFile(path, []byte("tag: [invalid"), 0644)

	cfg, err := pkg.LoadPackageConfig(path)
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to parse config")
}

func TestLoadPackageConfig_ValidYAML(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "valid.yaml")
	content := []byte(`
repository: github.com/core/project
tag: v1.2.3
`)
	os.WriteFile(path, content, 0644)

	cfg, err := pkg.LoadPackageConfig(path)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, path, cfg.Path)
	assert.Equal(t, "github.com/core/project", cfg.Repository)
	assert.Equal(t, "v1.2.3", cfg.Tag)
}
