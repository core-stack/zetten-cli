package auth

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/core-stack/zetten-cli/internal/util"
	"github.com/goccy/go-yaml"
)

var AUTH_FILE_NAME = "zetten-auth.yml"

type AuthConfig struct {
	Method      string `yaml:"method"`
	Credentials string `yaml:"credentials"`
}

type AuthMap map[string]AuthConfig

type AuthConfigLoader struct {
	configs AuthMap
}

func (l *AuthConfigLoader) FindAuth(repoUrl string) (*AuthConfig, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repo URL: %w", err)
	}

	host := u.Host
	urlPaths := strings.Split(u.Path, "/")

	// verify decremental repository path
	for i := len(urlPaths); i >= 1; i-- {
		path := urlPaths[:i]
		cfg, find := l.loadByHostPath(host, strings.Join(path, "/"))
		if find {
			return cfg, nil
		}
	}
	return nil, errors.New("Auth config not found for this repo")
}

func (l *AuthConfigLoader) loadByHostPath(host string, path string) (*AuthConfig, bool) {
	hostPath := fmt.Sprintf("%s/%s", host, path)

	// find config by host path
	cfg, exists := util.FindInMap(l.configs,
		func(k string, v AuthConfig) bool {
			return k == hostPath
		},
	)
	if exists {
		return &cfg, true
	} else {
		return nil, false
	}
}

func newAuthConfigLoader() *AuthConfigLoader {
	paths := []string{
		AUTH_FILE_NAME, // local
		filepath.Join(os.Getenv("HOME"), fmt.Sprintf(".zetten/%s", AUTH_FILE_NAME)), // global
	}
	configs := make(AuthMap)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var auths AuthMap
		if err := yaml.Unmarshal(data, &auths); err != nil {
			fmt.Println(fmt.Sprintf("invalid YAML in %s: %w", path, err))
		}

		configs = util.MergeMap(configs, auths)
	}

	return &AuthConfigLoader{
		configs: configs,
	}
}

var Loader = newAuthConfigLoader()
