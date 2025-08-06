package root_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/core-stack/zetten-cli/internal/core/root"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func readRootFile(path string) (*root.RootFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rf root.RootFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, err
	}
	return &rf, nil
}

func TestRootFile_Save(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "root.yaml")

	r := &root.RootFile{
		ZettenProjects: []string{"/home/user/project1"},
		Path:           path,
	}

	err := r.Save()
	assert.NoError(t, err)

	loaded, err := readRootFile(path)
	assert.NoError(t, err)
	assert.Equal(t, []string{"/home/user/project1"}, loaded.ZettenProjects)
}

func TestRootFile_AddProject(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "root.yaml")

	r := &root.RootFile{
		Path: path,
	}

	err := r.AddProject("/projects/z1", true)
	assert.NoError(t, err)
	assert.Contains(t, r.ZettenProjects, "/projects/z1")

	loaded, err := readRootFile(path)
	assert.NoError(t, err)
	assert.Contains(t, loaded.ZettenProjects, "/projects/z1")
}

func TestRootFile_RemoveProject(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "root.yaml")

	r := &root.RootFile{
		ZettenProjects: []string{"/projects/a", "/projects/b"},
		Path:           path,
	}
	r.Save()

	err := r.RemoveProject("/projects/a", true)
	assert.NoError(t, err)
	assert.NotContains(t, r.ZettenProjects, "/projects/a")
	assert.Contains(t, r.ZettenProjects, "/projects/b")

	loaded, err := readRootFile(path)
	assert.NoError(t, err)
	assert.NotContains(t, loaded.ZettenProjects, "/projects/a")
	assert.Contains(t, loaded.ZettenProjects, "/projects/b")
}

func TestRootFile_RemoveNonexistentProject(t *testing.T) {
	r := &root.RootFile{
		ZettenProjects: []string{"/x", "/y"},
	}

	err := r.RemoveProject("/not-found", false)
	assert.NoError(t, err)
	assert.Equal(t, []string{"/x", "/y"}, r.ZettenProjects)
}
