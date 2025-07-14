package main

import (
	"github.com/alecthomas/kong"
	"github.com/core-stack/zetten-cli/internal/commands/create"
	"github.com/core-stack/zetten-cli/internal/commands/initialize"
)

var cli struct {
	Init   initialize.InitCommand `cmd:"" help:"Initialize a new project."`
	Create create.CreateCommand   `cmd:"" help:"Create a new package."`
}

func main() {
	ctx := kong.Parse(&cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
