package promote

import (
	"github.com/core-stack/zetten-cli/internal/cli/prompt"
	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/core-stack/zetten-cli/internal/util"
)

type PromoteCommand struct {
	Url string `help:"The URL of the package to install" short:"u" long:"url"`
	Tag string `help:"The tag/version to install" short:"t" long:"tag"`

	config *project.ProjectConfig
}

func (c *PromoteCommand) BeforeApply() error {
	config, err := project.LoadProjectConfig("zetten.yml")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *PromoteCommand) Run() error {
	var err error
	if len(c.Url) == 0 {
		keys := util.MapKeys[map[string]string](c.config.Dependencies)
		c.Url, err = prompt.PromptSelect("Select a package to promote", keys, false)
		if err != nil {
			return err
		}
	}
	if c.Tag == "" {
		tag, err := prompt.PromptInput("Tag")
		if err != nil {
			return err
		}
		c.Tag = tag
	}

	return c.config.Promote(c.Url, c.Tag)
}
