package project_test

import (
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
)

// MockRootConfig finge o comportamento real
type MockRootConfig struct{}

func (m *MockRootConfig) OpenOrClonePackage(url string) (*git.Repository, error) {
	return nil, nil
}
func (m *MockRootConfig) Checkout(url, tag string) (*git.Repository, error) {
	if tag == "error" {
		return nil, errors.New("checkout failed")
	}
	return nil, nil
}
func (m *MockRootConfig) CopyRootFiles(url, destination string, ignore []string) error {
	return os.MkdirAll(destination, 0755)
}
func (m *MockRootConfig) Promote(url, tag, newTag, packageDir string) error {
	return nil
}
func (m *MockRootConfig) HasPackage(url string) bool {
	return true
}

func (m *MockRootConfig) BuildRootPackagePath(url string) string {
	return "/fake/path"
}
