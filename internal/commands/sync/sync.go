package sync

import "github.com/core-stack/zetten-cli/internal/core/project"

type SyncCommand struct {
	config *project.ProjectConfig
}

func (c *SyncCommand) BeforeApply() error {
	config, err := project.LoadProjectConfig("zetten.yml")
	if err != nil {
		return err
	}
	c.config = config
	return nil
}

func (c *SyncCommand) Run() error {
	return c.config.Sync()
}
