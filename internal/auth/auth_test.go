package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFiles(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "auth-test-*.yml")
	require.NoError(t, err)
	defer tmpfile.Close()

	_, err = tmpfile.Write([]byte(content))
	require.NoError(t, err)

	return tmpfile.Name()
}

func TestNewAuthConfigLoader(t *testing.T) {
	t.Run("with valid YAML", func(t *testing.T) {
		validYAML := `github.com/user/repo:
  method: token
  credentials: my-token
gitlab.com:
  method: oauth
  credentials: my-oauth-token`

		path := setupTestFiles(t, validYAML)
		defer os.Remove(path)

		loader := NewAuthConfigLoader(path)
		assert.NotNil(t, loader)
		assert.Equal(t, 2, len(loader.Configs))
	})

	t.Run("with invalid YAML", func(t *testing.T) {
		invalidYAML := `invalid: yaml: structure
github.com/user/repo:
  method: token
  credentials without colon`

		path := setupTestFiles(t, invalidYAML)
		defer os.Remove(path)

		loader := NewAuthConfigLoader(path)
		assert.NotNil(t, loader)
		assert.Empty(t, loader.Configs)
	})

	t.Run("with non-existent file", func(t *testing.T) {
		loader := NewAuthConfigLoader("non-existent-file.yml")
		assert.NotNil(t, loader)
		assert.Empty(t, loader.Configs)
	})

	t.Run("with empty filename", func(t *testing.T) {
		// This will try to load DEFAULT_AUTH_FILE_NAME which shouldn't exist in test env
		loader := NewAuthConfigLoader("")
		assert.NotNil(t, loader)
		assert.Empty(t, loader.Configs)
	})
}

func TestFindAuth(t *testing.T) {
	loader := &AuthConfigLoader{
		Configs: AuthMap{
			"github.com/user/repo": {
				Method:      "token",
				Credentials: "repo-token",
			},
			"github.com/user": {
				Method:      "oauth",
				Credentials: "user-token",
			},
			"gitlab.com": {
				Method:      "basic",
				Credentials: "gitlab-cred",
			},
		},
	}

	tests := []struct {
		name     string
		repoUrl  string
		expected *AuthConfig
		wantErr  bool
	}{
		{
			name:     "exact match",
			repoUrl:  "https://github.com/user/repo",
			expected: &AuthConfig{Method: "token", Credentials: "repo-token"},
		},
		{
			name:     "parent match",
			repoUrl:  "https://github.com/user/another-repo",
			expected: &AuthConfig{Method: "oauth", Credentials: "user-token"},
		},
		{
			name:     "host match",
			repoUrl:  "https://gitlab.com/group/project",
			expected: &AuthConfig{Method: "basic", Credentials: "gitlab-cred"},
		},
		{
			name:    "no match",
			repoUrl: "https://bitbucket.org/user/repo",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			repoUrl: "://invalid-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := loader.FindAuth(tt.repoUrl)
			fmt.Println(tt.repoUrl, cfg, err)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}

func TestLoadByHostPath(t *testing.T) {
	loader := &AuthConfigLoader{
		Configs: AuthMap{
			"github.com/user/repo": {
				Method:      "token",
				Credentials: "repo-token",
			},
		},
	}

	t.Run("existing path", func(t *testing.T) {
		cfg, found := loader.LoadByHostPath("github.com", "user/repo")
		assert.True(t, found)
		assert.Equal(t, "token", cfg.Method)
		assert.Equal(t, "repo-token", cfg.Credentials)
	})

	t.Run("non-existing path", func(t *testing.T) {
		cfg, found := loader.LoadByHostPath("github.com", "user/other-repo")
		assert.False(t, found)
		assert.Nil(t, cfg)
	})

	t.Run("empty path", func(t *testing.T) {
		cfg, found := loader.LoadByHostPath("github.com", "")
		assert.False(t, found)
		assert.Nil(t, cfg)
	})
}

func TestDefaultLoader(t *testing.T) {
	t.Run("default loader initialization", func(t *testing.T) {
		assert.NotNil(t, Loader)
		assert.IsType(t, &AuthConfigLoader{}, Loader)
	})
}

func TestMergeConfigs(t *testing.T) {
	t.Run("merge local and global configs", func(t *testing.T) {
		localYAML := `github.com/user/repo:
  method: token
  credentials: local-token`

		globalYAML := `github.com/user:
  method: oauth
  credentials: global-token
gitlab.com:
  method: basic
  credentials: global-cred`

		// Setup local file
		localPath := setupTestFiles(t, localYAML)
		defer os.Remove(localPath)

		// Setup global file in temp home dir
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		globalDir := filepath.Join(tmpHome, ".zetten")
		os.Mkdir(globalDir, 0755)
		globalPath := filepath.Join(globalDir, DEFAULT_AUTH_FILE_NAME)
		os.WriteFile(globalPath, []byte(globalYAML), 0644)

		loader := NewAuthConfigLoader(localPath)
		assert.NotNil(t, loader)
		assert.Equal(t, 3, len(loader.Configs))

		// Verify local config takes precedence
		cfg, found := loader.LoadByHostPath("github.com", "user/repo")
		assert.True(t, found)
		assert.Equal(t, "token", cfg.Method)
		assert.Equal(t, "local-token", cfg.Credentials)
	})
}
