package git_util

import (
	"fmt"
	"os"
	"strings"

	"github.com/core-stack/zetten-cli/internal/auth"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
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

func CloneRepo(repoUrl, destination string, opts ...CloneOpt) error {
	options := &CloneOptions{
		RepoUrl:     repoUrl,
		Destination: destination,
	}
	for _, opt := range opts {
		opt(options)
	}

	if err := os.MkdirAll(options.Destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	var authMethod transport.AuthMethod
	var err error

	if options.AuthMethod == "" || options.Credentials == "" {
		authConf, err := auth.Loader.FindAuth(options.RepoUrl)
		if err != nil {
			if err != auth.ErrNoAuthConfigFound {
				return fmt.Errorf("failed to load auth for url %s: %w", options.RepoUrl, err)
			} else {
				options.AuthMethod = "none"
				fmt.Printf("no auth found for %s, cloning without authentication\n", options.RepoUrl)
			}
		} else {
			options.AuthMethod = authConf.Method
			options.Credentials = authConf.Credentials
		}
	}

	switch options.AuthMethod {
	case "token":
		authMethod = &http.BasicAuth{
			Username: "git",
			Password: options.Credentials,
		}
	case "basic":
		parts := strings.Split(options.Credentials, ":")
		if len(parts) != 2 {
			return fmt.Errorf("basic auth credentials must be in format 'username:password'")
		}
		authMethod = &http.BasicAuth{
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

		authMethod = publicKeys
	case "none":
		authMethod = nil
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

	exists, err := RemoteRefExists(options.RepoUrl, referenceName, authMethod)
	if err != nil {
		return fmt.Errorf("failed to check remote ref: %w", err)
	}
	if !exists {
		return fmt.Errorf("reference %s does not exist in remote repository", referenceName)
	}

	repo, err := git.PlainClone(options.Destination, false, &git.CloneOptions{
		URL:           options.RepoUrl,
		Auth:          authMethod,
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

func ExtractBranchs(refs storer.ReferenceIter) []string {
	var branches []string
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, strings.TrimPrefix(ref.Name().String(), "refs/heads/"))
		}
		return nil
	})
	return branches
}
func ExtractTags(refs storer.ReferenceIter) []string {
	var tags []string
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsTag() {
			tags = append(tags, strings.TrimPrefix(ref.Name().String(), "refs/tags/"))
		}
		return nil
	})
	return tags
}
