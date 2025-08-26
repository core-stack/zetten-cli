package auth

import (
	"fmt"
	"path/filepath"

	"github.com/core-stack/zetten-cli/internal/util"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var DEFAULT_AUTH_FILE_NAME = "zetten-auth.yml"

type Scope string

const (
	ScopeLocal  Scope = "local"
	ScopeGlobal Scope = "global"
)

type Auth struct {
	Method      string `yaml:"method"`
	Credentials string `yaml:"credentials"`
}

func (a *Auth) ToAuthMethod() (transport.AuthMethod, error) {
	if a.Method == "ssh" {
		return ssh.DefaultAuthBuilder(a.Credentials)
	}
	if a.Method == "http" {
		return &http.BasicAuth{
			Username: "zetten",
			Password: a.Credentials,
		}, nil
	}
	return nil, fmt.Errorf("unknown auth method: %s", a.Method)
}

type AuthFile struct {
	Auths map[string]Auth `yaml:"auths"`
	Path  string          `yaml:"-"`
	Scope Scope           `yaml:"-"`
}

func (f *AuthFile) Save() error {
	configPath := filepath.Join(f.Path)
	return util.SaveYAMLIndented(configPath, f)
}

func (f *AuthFile) AddAuth(name string, method string, credentials string, autoSave bool) error {
	f.Auths[name] = Auth{Method: method, Credentials: credentials}
	if autoSave {
		return f.Save()
	}
	return nil
}

func (f *AuthFile) RemoveAuth(name string, autoSave bool) error {
	delete(f.Auths, name)
	if autoSave {
		return f.Save()
	}
	return nil
}
