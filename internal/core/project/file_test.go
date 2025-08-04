package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func readProjectFile(path string) (*project.ProjectFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf project.ProjectFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return &pf, nil
}

func TestProjectFile_Save(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")

	p := &project.ProjectFile{
		Name:         "my-project",
		Version:      "1.0.0",
		PackagesPath: "vendor",
		Dependencies: map[string]string{"pkg-a": "1.2.3"},
		Path:         path,
	}

	err := p.Save()
	assert.NoError(t, err)

	loaded, err := readProjectFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "my-project", loaded.Name)
	assert.Equal(t, "1.0.0", loaded.Version)
	assert.Equal(t, "vendor", loaded.PackagesPath)
	assert.Equal(t, "1.2.3", loaded.Dependencies["pkg-a"])
}

func TestProjectFile_AddDependency(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")

	p := &project.ProjectFile{
		Name:         "test",
		Version:      "0.0.1",
		Dependencies: map[string]string{},
		Path:         path,
	}

	err := p.AddDependency("example-lib", "2.1.0", true)
	assert.NoError(t, err)
	assert.Equal(t, "2.1.0", p.Dependencies["example-lib"])

	loaded, err := readProjectFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "2.1.0", loaded.Dependencies["example-lib"])
}

func TestProjectFile_RemoveDependency(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")

	p := &project.ProjectFile{
		Dependencies: map[string]string{
			"libx": "3.3.3",
			"liby": "4.4.4",
		},
		Path: path,
	}

	err := p.RemoveDependency("libx", true)
	assert.NoError(t, err)
	_, exists := p.Dependencies["libx"]
	assert.False(t, exists)

	loaded, err := readProjectFile(path)
	assert.NoError(t, err)
	_, exists = loaded.Dependencies["libx"]
	assert.False(t, exists)
	assert.Equal(t, "4.4.4", loaded.Dependencies["liby"])
}

func TestProjectFile_SetVersion(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "project.yaml")

	p := &project.ProjectFile{
		Version: "0.1.0",
		Path:    path,
	}

	err := p.SetVersion("0.2.0", true)
	assert.NoError(t, err)
	assert.Equal(t, "0.2.0", p.Version)

	loaded, err := readProjectFile(path)
	assert.NoError(t, err)
	assert.Equal(t, "0.2.0", loaded.Version)
}
