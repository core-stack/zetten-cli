package auth

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/core-stack/zetten-cli/config"
	"github.com/core-stack/zetten-cli/internal/core/file"
	"github.com/core-stack/zetten-cli/internal/util"
)

var (
	ErrNoAuthConfigFound = errors.New("no auth config found")
)

type IAuthConfig interface {
	FindAuth(repoUrl string) (*Auth, error)
	AddAuth(scope Scope, name string, method string, credentials string, autoSave bool) error
}

type AuthConfig struct {
	Global *AuthFile
	Local  *AuthFile
}

func (l *AuthConfig) FindAuth(repoUrl string) (*Auth, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repo URL: %w", err)
	}

	host := u.Host
	urlPaths := strings.Split(u.Path, "/")
	// verify decremental repository path
	for i := len(urlPaths); i >= 1; i-- {
		path := urlPaths[:i]
		cfg, found := l.loadByHostPath(host, strings.Join(path, "/"))
		if found {
			return cfg, nil
		}
	}
	defaultAuth := l.Global.Auths["default"]
	if defaultAuth.Credentials != "" {
		return &defaultAuth, nil
	}
	return nil, ErrNoAuthConfigFound
}

func (l *AuthConfig) loadByHostPath(host string, path string) (*Auth, bool) {
	hostPath := fmt.Sprintf("%s%s", host, path)
	fmt.Println(hostPath)
	// find config by host path
	if l.Local != nil {
		cfg, exists := util.FindInMap(l.Local.Auths,
			func(k string, v Auth) bool {
				return k == hostPath
			},
		)
		if exists {
			return &cfg, true
		}
	}
	cfg, exists := util.FindInMap(l.Global.Auths,
		func(k string, v Auth) bool {
			return k == hostPath
		},
	)
	if exists {
		return &cfg, true
	}
	return nil, false
}

func (l *AuthConfig) AddAuth(scope Scope, name string, method string, credentials string, autoSave bool) error {
	switch scope {
	case ScopeLocal:
		return l.Local.AddAuth(name, method, credentials, autoSave)
	case ScopeGlobal:
		return l.Global.AddAuth(name, method, credentials, autoSave)
	}
	return nil
}

type LoadAuthConfigsParams struct {
	Path string
}

func WithPath(path string) func(config *LoadAuthConfigsParams) {
	return func(config *LoadAuthConfigsParams) {
		config.Path = path
	}
}
func LoadAuthConfigs(opts ...func(config *LoadAuthConfigsParams)) (IAuthConfig, error) {
	var params LoadAuthConfigsParams
	for _, opt := range opts {
		opt(&params)
	}
	// check if root folder exists
	if _, err := os.Stat(config.DEFAULT_ROOT_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(config.DEFAULT_ROOT_PATH, os.ModePerm); err != nil {
			return nil, fmt.Errorf("❌ Failed to create root config directory: %w", err)
		}
	}

	// check if auth file exists
	if _, err := os.Stat(config.DEFAULT_ROOT_AUTHS_PATH); os.IsNotExist(err) {
		authFile := &AuthFile{
			Path:  config.DEFAULT_ROOT_AUTHS_PATH,
			Scope: ScopeGlobal,
		}
		if err = authFile.Save(); err != nil {
			return nil, fmt.Errorf("❌ Failed to create default auth file: %w", err)
		}
	}

	authConfig := AuthConfig{}

	// load global auth configs
	global, err := file.Load[AuthFile](config.DEFAULT_ROOT_AUTHS_PATH)
	if err != nil {
		return nil, err
	}
	global.Path = config.DEFAULT_ROOT_AUTHS_PATH
	global.Scope = ScopeGlobal
	authConfig.Global = global

	// load local auth configs
	if params.Path != "" {
		if _, err := os.Stat(params.Path); os.IsNotExist(err) {
			authFile := &AuthFile{
				Path:  params.Path,
				Scope: ScopeLocal,
			}
			if err = authFile.Save(); err != nil {
				return nil, fmt.Errorf("❌ Failed to create auth file: %w", err)
			}
		}
		local, err := file.Load[AuthFile](params.Path)
		if err != nil {
			return nil, err
		}
		local.Path = params.Path
		local.Scope = ScopeLocal
		authConfig.Local = local
	}

	return &authConfig, nil
}
