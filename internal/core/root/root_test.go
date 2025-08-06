package root_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/root"
	"github.com/stretchr/testify/assert"
)

func TestBuildRootPackagePath(t *testing.T) {
	r := &root.RootConfig{}
	url := "https://github.com/user/repo.git"
	expected := filepath.Join(root.DEFAULT_ROOT_PACKAGES_PATH, "user/repo")

	result := r.BuildRootPackagePath(url)
	assert.Equal(t, expected, result)
}

func TestHasPackage(t *testing.T) {
	r := &root.RootConfig{}

	// Deve retornar false
	assert.False(t, r.HasPackage("/non/existent/path"))

	// Criar pasta fake
	fakeURL := "https://example.com/org/proj.git"
	path := r.BuildRootPackagePath(fakeURL)
	os.MkdirAll(path, 0755)

	assert.True(t, r.HasPackage(fakeURL))
}

func TestLoadRootConfig_CreatesConfig(t *testing.T) {
	// Redefinindo paths globalmente (não ideal, mas necessário para testar isoladamente)
	originalPath := root.DEFAULT_ROOT_PATH
	root.DEFAULT_ROOT_PATH = t.TempDir()
	root.DEFAULT_ROOT_CONFIG_PATH = filepath.Join(root.DEFAULT_ROOT_PATH, "config.yml")
	root.DEFAULT_ROOT_PACKAGES_PATH = filepath.Join(root.DEFAULT_ROOT_PATH, "packages")
	defer func() { root.DEFAULT_ROOT_PATH = originalPath }()

	cfg, err := root.LoadRootConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.FileExists(t, root.DEFAULT_ROOT_CONFIG_PATH)
	assert.Equal(t, root.DEFAULT_ROOT_CONFIG_PATH, cfg.Path)
	assert.Empty(t, cfg.ZettenProjects)
}

func TestCopyRootFiles(t *testing.T) {
	r := &root.RootConfig{}
	tmpDst := t.TempDir()

	// Simula estrutura de pacote clonado
	url := "https://example.com/my/repo.git"
	srcPath := r.BuildRootPackagePath(url)
	os.MkdirAll(srcPath, 0755)
	os.WriteFile(filepath.Join(srcPath, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(srcPath, ".git"), []byte("ignore me"), 0644)

	err := r.CopyRootFiles(url, tmpDst, []string{})
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDst, "main.go"))
	assert.NoFileExists(t, filepath.Join(tmpDst, ".git"))
}

func TestOpenOrClonePackage(t *testing.T) {
	r := &root.RootConfig{}
	url := "https://github.com/octocat/Hello-World.git"

	repo, err := r.OpenOrClonePackage(url)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}
