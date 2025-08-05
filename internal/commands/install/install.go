package install

import (
	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/core-stack/zetten-cli/internal/git_util"
	"github.com/core-stack/zetten-cli/internal/prompt"
)

type InstallCommand struct {
	Url    string `help:"The URL of the package to install" short:"u" long:"url"`
	Tag    string `help:"The tag/version to install" short:"t" long:"tag"`
	Branch string `help:"The branch to install" short:"b" long:"branch"`

	config *project.ProjectConfig
}

func (c *InstallCommand) BeforeApply() error {
	config, err := project.LoadProjectConfig("zetten.yml")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *InstallCommand) Run() error {
	var err error
	if c.Url == "" {
		c.Url, err = prompt.PromptInput("📝 Package URL")
		if err != nil {
			return err
		}
	}
	repo, err := c.config.Root.OpenOrClonePackage(c.Url)
	if err != nil {
		return err
	}
	tagOrBranch, err := git_util.LoadBranchOrTag(repo, c.Branch, c.Tag)
	if err != nil {
		return err
	}
	c.config.Install(c.Url, tagOrBranch)
	return nil
}
