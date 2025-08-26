package root_test

import "github.com/core-stack/zetten-cli/internal/core/auth"

type MockAuthConfig struct {
	FindAuthFunc func(repoUrl string) (*auth.Auth, error)
}

func (m *MockAuthConfig) FindAuth(repoUrl string) (*auth.Auth, error) {
	return m.FindAuthFunc(repoUrl)
}

func (l *MockAuthConfig) AddAuth(scope auth.Scope, name string, method string, credentials string, autoSave bool) error {
	return nil
}
func NewMockAuthConfig() *MockAuthConfig {
	return &MockAuthConfig{
		FindAuthFunc: func(repoUrl string) (*auth.Auth, error) {
			return nil, nil
		},
	}
}
