package initialize

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/core-stack/zetten-cli/config"
	"github.com/core-stack/zetten-cli/internal/prompt"
	"github.com/core-stack/zetten-cli/internal/util"
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
	if c.Default {
		c.Name = getDefaultProjectName()
		c.Version = "1.0.0"
		c.Path = "packages"
	} else {
		var error error
		if c.Name == "" {
			c.Name, error = prompt.PromptInput("📝 Project name", getDefaultProjectName())
		}
		if error != nil {
			return error
		}

		if c.Version == "" {
			c.Version, error = prompt.PromptInput("📝 Project version", "1.0.0")
		}
		if error != nil {
			return error
		}

		if c.Path == "" {
			c.Path, error = prompt.PromptInput("📝 Packages path", "packages")
		}
		if error != nil {
			return error
		}
	}

	if _, err := os.Stat("zetten.yml"); err == nil {
		return errors.New("⚠️ A configuration file already exists at this location: " + "zetten.yml")
	}

	cfg := config.ProjectConfig{
		Name:    c.Name,
		Version: c.Version,
	}
	err := util.SaveYAMLIndented("zetten.yml", cfg)
	if err != nil {
		return fmt.Errorf("❌ Failed to save config: %w", err)
	}

	fmt.Println("✅ Project initialized successfully!")
	fmt.Println("📄 Config saved at:", "zetten.yml")
	return nil
}
