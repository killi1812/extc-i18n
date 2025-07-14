// Package translate provides main functionality extracting translate keys
package translate

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/subcommands"
	"go.uber.org/zap"
)

//nolint:allcaps
const (
	_COMMAND   = "grep"
	_CMD_ARGS  = "-rnbIE"
	_KEY_REGEX = "'+[A-Z]+\\.[A-Z_]+'"
)

type TranslateCmd struct {
	// projectPath is path to vue project
	projectPath string
	// outPath is path to translate folder
	outPath string
	// results stores locations of translate keys
	results result
}

// Execute implements subcommands.Command.
func (t *TranslateCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) subcommands.ExitStatus {
	zap.S().Debugf("Translate args: %+v", t)

	if err := t.Search(); err != nil {
		// Standardize errors
		return 1
	}
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

// result defines location of translate key
type resultold struct {
	// path to folder
	path string
	// line line of code
	line int
	// position refers to starting position in line
	position int

	//
	value keyValue
}

type keyValue struct {
	group string
	name  string
}

// Search searches for keys in project direcotry and fills results
func (t *TranslateCmd) Search() error {
	zap.S().Debugf("Searching project: %s", t.projectPath)

	ctx := context.TODO()
	cmd := exec.CommandContext(ctx, _COMMAND, _CMD_ARGS, _KEY_REGEX, t.projectPath)

	zap.S().Debugf("Executing cmd: %s with args: %s, with pattern: %s", _COMMAND, _CMD_ARGS, _KEY_REGEX)

	rez, err := cmd.Output()
	if err != nil {
		if exerr, ok := err.(*exec.ExitError); ok {
			zap.S().Debugf("Command %s exited with status code: %d", cmd.String(), exerr.ExitCode())

			// NOTE: grep exits with 1 if no results were selected and with 2 if error occurred
			if exerr.ExitCode() == 1 {
				zap.S().Infof("Search yielded no results")
			}
			if exerr.ExitCode() == 2 {
				zap.S().Errorf("Search (grep) failed with status code: %d, error: %s", exerr.ExitCode(), string(exerr.Stderr))
			}

		} else {
			zap.S().Debugf("Command: %s failed with err: %w", cmd.String(), err)
		}
		return err
	}

	if cmd.Err != nil {
		zap.S().Errorf("Command %s exited with err: %w", cmd.String(), cmd.Err)
		return cmd.Err
	}

	data := string(rez)
	zap.S().Debugf("Search resulted with %v", data)

	dataRez, err := Parse(data)
	if err != nil {
		zap.S().Errorf("Failed to Parse data error: %w", err)
	}

	t.results = dataRez
	fmt.Printf("t.results: %v\n", t.results)

	return nil
}

func Parse(data string) (result, error) {
	lines := strings.Split(data, "\n")
	lines = lines[:len(lines)-1]
	results := make(result)

	for i, line := range lines {
		zap.S().Debugf("Working on index: %d, line: %s", i, line)
		parts := strings.Split(line, ":")

		var arg args
		// Set path
		arg.path = strings.TrimSpace(parts[0])

		// Set line
		tmp, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			zap.S().DPanicf("Failed to parse result line got: %v, expected int", parts[1])
		}
		arg.line = tmp

		// Set position
		tmp, err = strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil {
			zap.S().DPanicf("Failed to parse result line got: %v, expected int", parts[2])
		}
		arg.position = tmp

		// Extracting key values
		{
			local := parts[3]
			zap.S().Debugf("Extracting key from %s", local)

			iprefix := indexRune(local, '\'')
			local = local[iprefix+1:]
			zap.S().Debugf("Removed left extra %s", local)

			isuffix := indexRune(local, '\'')
			local = local[:isuffix]
			zap.S().Debugf("Removed right extra %s", local)

			parts := strings.Split(local, ".")
			if len(parts) != 2 {
				zap.S().DPanicf("Bad input (%s), needs look like: GROUP.NAME", data)
			}

			group := parts[0]
			arg.name = parts[1]

			results.Add(group, arg)
		}
	}

	return results, nil
}

// indexRune will find first occurence of rune
//
// if there is no rune in given string it will return -1
func indexRune(data string, c rune) int {
	for i, char := range data {
		if char == c {
			return i
		}
	}
	return -1
}
