package promote

import (
	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/core-stack/zetten-cli/internal/git_util"
	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/core-stack/zetten-cli/internal/util"
)

type PromoteCommand struct {
	Url    string `help:"The URL of the package to install" short:"u" long:"url"`
	Tag    string `help:"The tag/version to install" short:"t" long:"tag"`
	Branch string `help:"The branch to install" short:"b" long:"branch"`

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
