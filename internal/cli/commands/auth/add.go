package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/core-stack/zetten-cli/internal/cli/prompt"
	auth_core "github.com/core-stack/zetten-cli/internal/core/auth"
)

type AddCommand struct {
	Local      bool   `short:"l" long:"local" help:"Add authentication credentials to the local config"`
	Type       string `short:"t" long:"type" help:"Authentication type"`
	Path       string `short:"p" long:"path" help:"Authentication path" default:"default"`
	Credential string `short:"c"  long:"credential" help:"Authentication credentials"`

	auth auth_core.IAuthConfig
}

func (c *AddCommand) BeforeApply() error {
	var err error
	if c.Local {
		c.auth, err = auth_core.LoadAuthConfigs(auth_core.WithPath("./zetten-auth.yml"))
	} else {
		c.auth, err = auth_core.LoadAuthConfigs()
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *AddCommand) Run() error {
	var err error
	if c.Type == "" {
		c.Type, err = prompt.PromptSelect("Select authentication type", []string{"http", "ssh"}, false)
	}
	if err != nil {
		return err
	}

	switch c.Type {
	case "http":
		err = c.handleHTTP()
	case "ssh":
		err = c.handleSSH()
	default:
		return fmt.Errorf("invalid authentication type: %s", c.Type)
	}
	if err != nil {
		return err
	}
	scope := auth_core.ScopeGlobal
	if c.Local {
		scope = auth_core.ScopeLocal
	}
	return c.auth.AddAuth(scope, c.Path, c.Type, c.Credential, true)
}

func (c *AddCommand) handleHTTP() error {
	var err error

	if c.Credential == "" {
		c.Credential, err = prompt.PromptInput("Token")
		if err != nil {
			return err
		}
	}

	if strings.TrimSpace(c.Credential) == "" {
		return fmt.Errorf("token cannot be empty")
	}

	return nil
}

func (c *AddCommand) handleSSH() error {
	if c.Credential == "" {
		sshKeys, err := c.findSSHKeys()
		if err != nil {
			return err
		}

		if len(sshKeys) == 0 {
			return fmt.Errorf("no SSH keys found in ~/.ssh/")
		}

		options := make([]string, len(sshKeys))
		for i, key := range sshKeys {
			options[i] = key.Path
		}

		options = append(options, "Specify another path...")

		selected, err := prompt.PromptSelect("Select a SSH key", options, false)
		if err != nil {
			return err
		}

		if selected == "Specify another path..." {
			c.Credential, err = prompt.PromptInput("Path to SSH")
			if err != nil {
				return err
			}
		} else {
			c.Credential = selected
		}
	}

	if _, err := os.Stat(c.Credential); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", c.Credential)
	}

	return nil
}

type SSHKey struct {
	Path    string
	Name    string
	Type    string
	IsValid bool
}

func (c *AddCommand) findSSHKeys() ([]SSHKey, error) {
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	var keys []SSHKey

	files, err := os.ReadDir(sshDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read .ssh: %v", err)
	}

	privateKeyExts := []string{"", ".pub", ".ppk"}
	keyTypes := []string{"rsa", "dsa", "ecdsa", "ed25519"}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		filePath := filepath.Join(sshDir, filename)

		if strings.HasSuffix(filename, ".pub") ||
			filename == "known_hosts" ||
			filename == "config" ||
			strings.HasPrefix(filename, ".") {
			continue
		}

		for _, ext := range privateKeyExts {
			for _, keyType := range keyTypes {
				patterns := []string{
					"id_" + keyType + ext,
					"id_" + keyType,
					keyType + "_key",
				}

				for _, pattern := range patterns {
					if filename == pattern || strings.Contains(filename, keyType) {
						keys = append(keys, SSHKey{
							Path:    filePath,
							Name:    filename,
							Type:    keyType,
							IsValid: true,
						})
						break
					}
				}
			}
		}
	}

	if len(keys) == 0 {
		for _, file := range files {
			if file.IsDir() || strings.HasSuffix(file.Name(), ".pub") {
				continue
			}

			filePath := filepath.Join(sshDir, file.Name())
			keys = append(keys, SSHKey{
				Path:    filePath,
				Name:    file.Name(),
				Type:    "unknown",
				IsValid: true,
			})
		}
	}

	return keys, nil
}
