package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/file"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

func TestLoad_EmptyPath(t *testing.T) {
	cfg, err := file.Load[TestConfig]("")
	assert.Error(t, err)
	assert.Equal(t, "", cfg.Name)
}

func TestLoad_FileNotFound(t *testing.T) {
	cfg, err := file.Load[TestConfig]("nonexistent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Configuration file not found")
	assert.Equal(t, "", cfg.Name)
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "invalid.yaml")
	os.WriteFile(path, []byte("name: João\nvalue: [invalid"), 0644)

	cfg, err := file.Load[TestConfig](path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to parse config")
	assert.Equal(t, "", cfg.Name)
}

func TestLoad_ReadError(t *testing.T) {
	// Cria arquivo e remove permissão de leitura
	tmp := t.TempDir()
	path := filepath.Join(tmp, "unreadable.yaml")
	os.WriteFile(path, []byte("name: test\nvalue: 42"), 0000)
	defer os.Chmod(path, 0644) // restaura para não travar cleanup

	cfg, err := file.Load[TestConfig](path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Failed to read config")
	assert.Equal(t, "", cfg.Name)
}

func TestLoad_ValidYAML(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "valid.yaml")
	content := []byte("name: testuser\nvalue: 42")
	os.WriteFile(path, content, 0644)

	cfg, err := file.Load[TestConfig](path)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", cfg.Name)
	assert.Equal(t, 42, cfg.Value)
}
