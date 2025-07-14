package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// CloneRepo clona um repositório em uma tag específica
// repoURL: URL do repositório (https/ssh)
// destination: pasta de destino
// tag: tag a ser clonada
// authMethod: "none", "token", "ssh" ou "basic"
// credentials: token, senha ou caminho para chave SSH
func CloneRepo(repoURL, destination, tag, authMethod, credentials string) error {
	// Verificar se o destino existe ou criá-lo
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Configurar autenticação
	var auth transport.AuthMethod
	var err error

	switch authMethod {
	case "token":
		// Para GitHub, GitLab com token
		auth = &http.BasicAuth{
			Username: "git", // Pode ser qualquer coisa para token auth
			Password: credentials,
		}
	case "basic":
		// Para Bitbucket ou autenticação básica
		parts := strings.Split(credentials, ":")
		if len(parts) != 2 {
			return fmt.Errorf("basic auth credentials must be in format 'username:password'")
		}
		auth = &http.BasicAuth{
			Username: parts[0],
			Password: parts[1],
		}
	case "ssh":
		// Autenticação por SSH
		publicKeys, err := ssh.NewPublicKeysFromFile("git", credentials, "")
		if err != nil {
			return fmt.Errorf("failed to parse SSH key: %w", err)
		}

		// Opcional: verificação de host conhecido
		hostKeyCallback, err := knownhosts.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
		if err == nil {
			publicKeys.HostKeyCallback = hostKeyCallback
		}

		auth = publicKeys
	case "none":
		// Sem autenticação (repositório público)
		auth = nil
	default:
		return fmt.Errorf("invalid auth method: %s", authMethod)
	}

	// Clonar o repositório
	repo, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:           repoURL,
		Auth:          auth,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tag)),
		SingleBranch:  true,
		Depth:         1,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Verificar se a tag foi realmente encontrada
	_, err = repo.Tag(tag)
	if err != nil {
		return fmt.Errorf("failed to find tag %s: %w", tag, err)
	}

	return nil
}
