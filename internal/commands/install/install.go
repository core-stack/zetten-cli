package install

import (
	"errors"
	"fmt"
	"strings"

	"github.com/core-stack/zetten-cli/config"
	"github.com/core-stack/zetten-cli/internal/git_util"
	"github.com/core-stack/zetten-cli/internal/util"
)

type InstallCommand struct {
	Url         string `help:"The URL of the package to install" short:"u" long:"url"`
	Provider    string `help:"The provider of the package to install" short:"p" long:"provider"`
	ProviderUrl string `help:"The provider url for custom git providers (eg: gitlab self hosted)" short:"purl" long:"provider_url"`
	Name        string `help:"The name of the package to install" short:"n" long:"name"`
	Tag         string `help:"The tag/version to install" short:"t" long:"tag"`
	Branch      string `help:"The branch to install" short:"b" long:"branch"`

	config *config.ProjectConfig
}

func (c *InstallCommand) BeforeApply() error {
	config, err := config.LoadProjectConfig(".")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *InstallCommand) Run() error {
	if c.Provider == "" {
		c.Provider = c.config.DefaultProvider
	}

	if c.Url == "" && (c.Provider == "" || c.Name == "") {
		return fmt.Errorf("you must provide either a full URL or both provider and name")
	}
	if c.Tag == "" && c.Branch == "" {
		return fmt.Errorf("you must provide either a tag or a branch")
	}

	repoURL, err := c.buildRepoURL()
	if err != nil {
		return err
	}

	destination, err := c.buildDestinationPath()
	if err != nil {
		return err
	}

	opts := []git_util.CloneOpt{}
	if c.Tag != "" {
		opts = append(opts, git_util.WithTag(c.Tag))
	} else if c.Branch != "" {
		opts = append(opts, git_util.WithBranch(c.Branch))
	}
	c.config.AddDependency(c.Url, util.Or(c.Tag, c.Branch))
	if err = c.config.Save(); err != nil {
		return errors.New("error saving new dependency")
	}

	return git_util.CloneRepo(repoURL, destination, opts...)
}

func (c *InstallCommand) buildRepoURL() (string, error) {
	if c.Url != "" {
		return c.normalizeURL(c.Url), nil
	}

	// check if is url
	if util.IsValidURL(c.Provider) {
		return fmt.Sprintf("%s/%s.git", c.Provider, c.Name), nil
	}

	// Construir URL baseada no provider
	switch strings.ToLower(c.Provider) {
	case "github":
		return fmt.Sprintf("https://github.com/%s.git", c.Name), nil
	case "gitlab":
		return fmt.Sprintf("https://gitlab.com/%s.git", c.Name), nil
	case "bitbucket":
		return fmt.Sprintf("https://bitbucket.org/%s.git", c.Name), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", c.Provider)
	}
}

func (c *InstallCommand) normalizeURL(url string) string {
	// Garantir que a URL termina com .git
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}

	// Converter URL SSH para HTTPS se necessário
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}

	return url
}

func (c *InstallCommand) buildDestinationPath() (string, error) {
	var repoName string

	if c.Name != "" {
		parts := strings.Split(c.Name, "/")
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid name format, should be 'owner/repo'")
		}
		repoName = parts[len(parts)-1]
	} else {
		// Extrair do URL
		parts := strings.Split(strings.TrimSuffix(c.Url, ".git"), "/")
		repoName = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s/%s", c.config.Path, repoName), nil
}
