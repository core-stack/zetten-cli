package pkg_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/pkg"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func readPackageFile(path string) (*pkg.PackageFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf pkg.PackageFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return &pf, nil
}

func TestPackageFile_GetName(t *testing.T) {
	p := &pkg.PackageFile{Repository: "github.com/user/repo.git"}
	assert.Equal(t, "repo", p.GetName())

	p.Repository = "gitlab.com/org/project"
	assert.Equal(t, "project", p.GetName())

	p.Repository = ""
	assert.Equal(t, "", p.GetName())
}

func TestPackageFile_SetTag(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "pkg.yaml")

	p := &pkg.PackageFile{Repository: "github.com/test/repo", Path: path}
	err := p.SetTag("v1.0.0", true)
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", p.Tag)

	loaded, err := readPackageFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "v1.0.0", loaded.Tag)
}

func TestPackageFile_SetRepository(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "pkg.yaml")

	p := &pkg.PackageFile{Path: path}
	err := p.SetRepository("github.com/test/project", true)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/test/project", p.Repository)

	loaded, err := readPackageFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "github.com/test/project", loaded.Repository)
}

func TestPackageFile_UpdateFrom(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "pkg.yaml")

	p := &pkg.PackageFile{
		Repository: "old.com/repo",
		Tag:        "v0.1.0",
		Path:       path,
	}
	p.Save()

	other := &pkg.PackageFile{
		Repository: "new.com/repo",
		Tag:        "1.0.0",
	}

	err := p.UpdateFrom(other, true)
	assert.NoError(t, err)
	assert.Equal(t, "new.com/repo", p.Repository)
	assert.Equal(t, "1.0.0", p.Tag)
	assert.Equal(t, "", p.Tag)

	loaded, err := readPackageFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "new.com/repo", loaded.Repository)
	assert.Equal(t, "1.0.0", loaded.Tag)
	assert.Equal(t, "", loaded.Tag)
}
