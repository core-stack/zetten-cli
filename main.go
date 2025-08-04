package main

import (
	"github.com/alecthomas/kong"
	"github.com/core-stack/zetten-cli/internal/commands/initialize"
	"github.com/core-stack/zetten-cli/internal/commands/install"
)

var cli struct {
	Init    initialize.InitCommand `cmd:"" help:"Initialize a new project."`
	Install install.InstallCommand `cmd:"" help:"Install a package."`
}

func main() {
	ctx := kong.Parse(&cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
