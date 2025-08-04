package install

import (
	"github.com/core-stack/zetten-cli/internal/core/project"
	"github.com/core-stack/zetten-cli/internal/git_util"
	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/core-stack/zetten-cli/internal/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	if c.Tag == "" && c.Branch == "" {
		if err := c.SelectTagOrBranch(repo); err != nil {
			return err
		}
	} else {
		if c.Tag != "" {
			_, err = repo.Tag(c.Tag)
			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					if err = c.SelectTagOrBranch(repo); err != nil {
						return err
					}
				}
				return err
			}
		} else if c.Branch != "" {
			_, err = repo.Branch(c.Branch)
			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					if err = c.SelectTagOrBranch(repo); err != nil {
						return err
					}
				}
				return err
			}
		}
	}
	c.config.Install(c.Url, util.Or(c.Tag, c.Branch))
	return nil
}

func (c *InstallCommand) SelectTag(repo *git.Repository) error {
	iterator, err := repo.Tags()
	if err != nil {
		return err
	}
	tags := git_util.ExtractTags(iterator)

	c.Tag, err = prompt.PromptSelect("📝 Tag", tags, true)
	if err != nil {
		return err
	}
	return nil
}

func (c *InstallCommand) SelectBranch(repo *git.Repository) error {
	branchs, err := repo.Branches()
	if err != nil {
		return err
	}
	c.Branch, err = prompt.PromptSelect("📝 Branch", git_util.ExtractBranchs(branchs), true)
	if err != nil {
		return err
	}
	return nil
}

func (c *InstallCommand) SelectTagOrBranch(repo *git.Repository) error {
	for c.Tag == "" && c.Branch == "" {
		selected, err := prompt.PromptSelect("📝 Select tag or branch", []string{"tag", "branch"}, false)
		if err != nil {
			return err
		}
		if selected == "tag" {
			err = c.SelectTag(repo)
		} else {
			err = c.SelectBranch(repo)
		}
		if err != nil && err != prompt.GoBack {
			return err
		}
	}
	return nil
}
