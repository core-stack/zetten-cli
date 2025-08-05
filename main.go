package main

import (
	"github.com/alecthomas/kong"
	"github.com/core-stack/zetten-cli/internal/commands/initialize"
	"github.com/core-stack/zetten-cli/internal/commands/install"
	"github.com/core-stack/zetten-cli/internal/commands/promote"
	"github.com/core-stack/zetten-cli/internal/commands/sync"
	"github.com/core-stack/zetten-cli/internal/commands/uninstall"
)

var cli struct {
	Init      initialize.InitCommand     `cmd:"" help:"Initialize a new project."`
	Install   install.InstallCommand     `cmd:"" help:"Install a package."`
	Uninstall uninstall.UninstallCommand `cmd:"" help:"Uninstall a package."`
	Update    install.InstallCommand     `cmd:"" help:"Update a package."`
	Sync      sync.SyncCommand           `cmd:"" help:"Sync packages."`
	Promote   promote.PromoteCommand     `cmd:"" help:"Promote a package."`
}

func main() {
	ctx := kong.Parse(&cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
