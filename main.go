package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/killi1812/extc-i18n/app"
	"github.com/killi1812/extc-i18n/cmd/translate"
)

func init() {
	app.Setup()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&translate.TranslateCmd{}, "")

	flag.Parse()
}

func main() {
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
