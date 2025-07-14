package create

import (
	"fmt"
	"path"

	"github.com/core-stack/zetten-cli/config"
	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/core-stack/zetten-cli/internal/util"
)

type CreateCommand struct {
	Name       string `help:"The name of the package" short:"n" long:"name"`
	Version    string `help:"The version of the package" short:"v" long:"version"`
	Private    bool   `help:"Create a private package" short:"p" long:"private"`
	Provider   string `help:"The provider of the package"` // github, gitlab, bitbucket
	Repository string `help:"The repository of the package"`

	config *config.ProjectConfig
}

func (c *CreateCommand) BeforeApply() error {
	config, err := config.LoadProjectConfig(".")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *CreateCommand) Run() error {
	var err error
	if c.Name == "" {
		c.Name, err = prompt.PromptInput("📝 Package name", "my-package")
	}
	if err != nil {
		return err
	}

	if c.Version == "" {
		c.Version, err = prompt.PromptInput("📝 Package version", "1.0.0")
	}
	if err != nil {
		return err
	}

	if !c.Private {
		c.Private, err = prompt.PromptConfirm("📝 Package is private", false)
	}
	if err != nil {
		return err
	}

	if c.Provider == "" {
		c.Provider, err = prompt.PromptInput("📝 Package provider", "github")
	}
	if err != nil {
		return err
	}

	if c.Repository == "" {
		c.Repository, err = prompt.PromptInput("📝 Package repository", "core-stack/zetten")
	}
	if err != nil {
		return err
	}

	configPath := path.Join(c.config.Path, c.Name, "zetten-config.yml")
	cfg := config.PackageConfig{
		Name:         c.Name,
		Version:      c.Version,
		Private:      c.Private,
		Provider:     c.Provider,
		Dependencies: config.Dependency{},
		Repository:   c.Repository,
	}
	err = util.SaveYAMLIndented(configPath, cfg)
	if err != nil {
		return fmt.Errorf("❌ Failed to save config: %w", err)
	}
	return nil
}
