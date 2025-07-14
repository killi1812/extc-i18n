// Package translate provides main functionality extracting translate keys
package translate

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

type TranslateCmd struct {
	// projectPath is path to vue project
	projectPath string

	// outPath is path to translate folder
	outPath string
}

// Execute implements subcommands.Command.
func (t *TranslateCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	zap.S().Debugf("cli args: %+v", t)
	return 0
}

// Name implements subcommands.Command.
func (t *TranslateCmd) Name() string { return "trn" }

// SetFlags implements subcommands.Command.
func (t *TranslateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&t.projectPath, "p", "./", "Path to your vue project")
	f.StringVar(&t.outPath, "o", "./", "Path to your language folder")
}

// Synopsis implements subcommands.Command.
func (t *TranslateCmd) Synopsis() string {
	return "extracts translate keys from vue project to translate folder"
}

// Usage implements subcommands.Command.
func (t *TranslateCmd) Usage() string {
	return `
	trn extracts translate keys from vue project to translate folder
	keys need to look like this to be translate 'PREFIX.KEY'
	trn -p /path/to/project -o /path/to/languages/directory\n`
}
