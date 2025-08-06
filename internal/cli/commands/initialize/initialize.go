package initialize

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/core-stack/zetten-cli/internal/cli/prompt"
	"github.com/core-stack/zetten-cli/internal/core/project"
)

type InitCommand struct {
	Name    string `help:"The name of the project" short:"n" long:"name"`
	Version string `help:"The version of the project" short:"v" long:"version"`
	Path    string `help:"The path to the packages" short:"p" long:"path"`

	Default bool `help:"Initialize a new project with default values." short:"d"`
}

func getDefaultProjectName() string {
	wd, err := os.Getwd()
	if err != nil {
		return "my-project" // fallback
	}
	return filepath.Base(wd)
}

func (c *InitCommand) Run() error {
	if _, err := os.Stat("zetten.yml"); err == nil {
		return errors.New("⚠️ A configuration file already exists at this location")
	}

	if c.Default {
		c.Name = getDefaultProjectName()
		c.Version = "1.0.0"
		c.Path = "packages"
	} else {
		var err error
		if c.Name == "" {
			c.Name, err = prompt.PromptInput("📝 Project name", prompt.WithDefaultValue(getDefaultProjectName()))
		}
		if err != nil {
			return err
		}

		if c.Version == "" {
			c.Version, err = prompt.PromptInput("📝 Project version", prompt.WithDefaultValue("1.0.0"))
		}
		if err != nil {
			return err
		}

		if c.Path == "" {
			c.Path, err = prompt.PromptInput("📝 Packages path", prompt.WithDefaultValue("packages"))
		}
		if err != nil {
			return err
		}
	}

	_, err := project.NewProjectConfig("zetten.yml", c.Name, c.Version, c.Path)

	if err != nil {
		return err
	}
	fmt.Println("✅ Project initialized successfully!")
	fmt.Println("📄 Config saved at:", "zetten.yml")
	return nil
}
