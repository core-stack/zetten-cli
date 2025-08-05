package uninstall

import (
	"errors"

	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/core-stack/zetten-cli/internal/util"
)

type UninstallCommand struct {
	Urls []string `help:"Comma-separated list of package URLs to uninstall" short:"u" long:"url" sep:","`

	config *project.ProjectConfig
}

func (c *UninstallCommand) BeforeApply() error {
	config, err := project.LoadProjectConfig("zetten.yml")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *UninstallCommand) Run() error {
	if c.config.Dependencies == nil || len(c.config.Dependencies) == 0 {
		return errors.New("no dependencies found")
	}
	if len(c.Urls) == 0 {
		keys := util.MapKeys[map[string]string](c.config.Dependencies)
		url, err := prompt.PromptSelect("Select packages to remove", keys, false)
		c.Urls = []string{url}
		if err != nil {
			return err
		}
	}
	return c.config.Uninstall(c.Urls)
}
