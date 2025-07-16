package util

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/core-stack/zetten-cli/config"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type CloneOptions struct {
	RepoUrl     string
	Destination string
	Tag         string
	Branch      string
	AuthMethod  string
	Credentials string
}

type CloneOpt func(*CloneOptions)

func WithRepoUrl(repoUrl string) CloneOpt {
	return func(o *CloneOptions) {
		o.RepoUrl = repoUrl
	}
}

func WithDestination(destination string) CloneOpt {
	return func(o *CloneOptions) {
		o.Destination = destination
	}
}

func WithTag(tag string) CloneOpt {
	return func(o *CloneOptions) {
		o.Tag = tag
	}
}

func WithBranch(branch string) CloneOpt {
	return func(o *CloneOptions) {
		o.Branch = branch
	}
}

func WithAuth(authMethod, credentials string) CloneOpt {
	return func(o *CloneOptions) {
		o.AuthMethod = authMethod
		o.Credentials = credentials
	}
}

func CloneRepo(opts ...CloneOpt) error {
	options := &CloneOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if err := os.MkdirAll(options.Destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	var auth transport.AuthMethod
	var err error
	u, err := url.Parse(options.RepoUrl)
	if err != nil {
		return fmt.Errorf("failed to parse repo URL: %w", err)
	}
	if options.AuthMethod == "" || options.Credentials == "" {
		authConf, err := config.LoadAuthForHost(u.Host)
		if err != nil {
			return fmt.Errorf("failed to load auth for host %s: %w", u.Host, err)
		}
		options.AuthMethod = authConf.Method
		options.Credentials = authConf.Credentials
	}

	switch options.AuthMethod {
	case "token":
		auth = &http.BasicAuth{
			Username: "git",
			Password: options.Credentials,
		}
	case "basic":
		parts := strings.Split(options.Credentials, ":")
		if len(parts) != 2 {
			return fmt.Errorf("basic auth credentials must be in format 'username:password'")
		}
		auth = &http.BasicAuth{
			Username: parts[0],
			Password: parts[1],
		}
	case "ssh":
		publicKeys, err := ssh.NewPublicKeysFromFile("git", options.Credentials, "")
		if err != nil {
			return fmt.Errorf("failed to parse SSH key: %w", err)
		}

		hostKeyCallback, err := knownhosts.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
		if err == nil {
			publicKeys.HostKeyCallback = hostKeyCallback
		}

		auth = publicKeys
	case "none":
		auth = nil
	default:
		return fmt.Errorf("invalid auth method: %s", options.AuthMethod)
	}

	var referenceName string
	if options.Tag != "" {
		referenceName = fmt.Sprintf("refs/tags/%s", options.Tag)
	} else if options.Branch != "" {
		referenceName = fmt.Sprintf("refs/heads/%s", options.Branch)
	} else {
		return fmt.Errorf("either tag or branch must be specified")
	}

	exists, err := RemoteRefExists(options.RepoUrl, referenceName, auth)
	if err != nil {
		return fmt.Errorf("failed to check remote ref: %w", err)
	}
	if !exists {
		return fmt.Errorf("reference %s does not exist in remote repository", referenceName)
	}

	repo, err := git.PlainClone(options.Destination, false, &git.CloneOptions{
		URL:           options.RepoUrl,
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName(referenceName),
		SingleBranch:  true,
		Depth:         1,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	if options.Tag != "" {
		_, err = repo.Tag(options.Tag)
		if err != nil {
			return fmt.Errorf("failed to find tag %s: %w", options.Tag, err)
		}
	}

	if options.Branch != "" {
		_, err = repo.Branch(options.Branch)
		if err != nil {
			return fmt.Errorf("failed to find branch %s: %w", options.Branch, err)
		}
	}

	return nil
}

func RemoteRefExists(repoURL, refTarget string, auth transport.AuthMethod) (bool, error) {
	remote := git.NewRemote(nil, &gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return false, fmt.Errorf("failed to list remote references: %w", err)
	}

	for _, ref := range refs {
		if ref.Name().String() == refTarget {
			return true, nil
		}
	}

	return false, nil
}
